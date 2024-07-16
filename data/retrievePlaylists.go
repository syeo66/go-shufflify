package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	. "github.com/syeo66/go-shufflify/types"
)

func RetrievePlaylists(token string, db *sql.DB) ([]Playlist, error) {
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
