package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "github.com/syeo66/go-shufflify/types"
)

func RetrieveFavouriteCount(token string, db *sql.DB) int {
	client := &http.Client{}
	requestURL := "https://api.spotify.com/v1/me/tracks?limit=1"
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	countObj := &SavedTracks{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, countObj)

	if err != nil {
		return 0
	}

	return countObj.Total
}
