package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	. "github.com/syeo66/go-shufflify/types"
)

func RetrieveNthSongFromFavourites(token string, n int) (*Track, error) {
	client := &http.Client{}
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?limit=1&offset=%d", n)
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving favourites"))
	}
	defer resp.Body.Close()

	tracks := &SavedTracks{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, tracks)

	if err != nil {
		return nil, errors.Join(err, errors.New("error unmarshalling favourites"))
	}

	if len(tracks.Items) == 0 {
		return nil, errors.New("no tracks found")
	}

	return &tracks.Items[0].Track, nil
}
