package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	. "github.com/syeo66/go-shufflify/types"
)

func RetrievePlaylists(token string, db *sql.DB) ([]Playlist, error) {
	_, found := CacheStore.Get(LOCK_KEY)
	if found {
		return nil, errors.New("request locked due to rate limiting")
	}

	client := &http.Client{}
	allPlaylists := []Playlist{}

	requestURL := "https://api.spotify.com/v1/me/playlists?limit=50"

	for requestURL != "" {
		req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := client.Do(req)
		if err != nil {
			return nil, errors.Join(err, errors.New("error retrieving playlists"))
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 {
			retryAfter := resp.Header.Get("Retry-After")
			r, err := strconv.Atoi(retryAfter)
			if err == nil {
				fmt.Printf("RetrievePlaylists rate limited, retry after: %s seconds\n", retryAfter)
				CacheStore.Set(LOCK_KEY, "true", time.Duration(r)*time.Second)
			}
			return nil, errors.New("rate limited by Spotify API")
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("Spotify API returned status %d when retrieving playlists", resp.StatusCode)
		}

		playlists := &Playlists{}
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, playlists)

		if err != nil {
			return nil, errors.Join(err, errors.New("error unmarshalling playlists"))
		}

		allPlaylists = append(allPlaylists, playlists.Items...)
		requestURL = playlists.Next
	}

	return allPlaylists, nil
}
