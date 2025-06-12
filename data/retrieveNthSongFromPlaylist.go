package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	. "github.com/syeo66/go-shufflify/types"
)

func RetrieveNthSongFromPlaylist(token string, p Playlist, n int) (*Track, error) {
	_, found := CacheStore.Get(LOCK_KEY)
	if found {
		return nil, errors.New("request locked due to rate limiting")
	}

	client := &http.Client{}
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?limit=1&offset=%d", p.Id, n)
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving playlist tracks"))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		retryAfter := resp.Header.Get("Retry-After")
		r, err := strconv.Atoi(retryAfter)
		if err == nil {
			fmt.Printf("RetrieveNthSongFromPlaylist rate limited, retry after: %s seconds\n", retryAfter)
			CacheStore.Set(LOCK_KEY, "true", time.Duration(r)*time.Second)
		}
		return nil, errors.New("rate limited by Spotify API")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Spotify API returned status %d when retrieving playlist tracks", resp.StatusCode)
	}

	tracks := &SavedTracks{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, tracks)

	if err != nil {
		return nil, errors.Join(err, errors.New("error unmarshalling playlist tracks"))
	}

	if len(tracks.Items) == 0 {
		return nil, errors.New("no tracks found")
	}

	return &tracks.Items[0].Track, nil
}
