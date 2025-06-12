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

func RetrieveFavouriteCount(token string, db *sql.DB) (int, error) {
	_, found := CacheStore.Get(LOCK_KEY)
	if found {
		return 0, errors.New("request locked due to rate limiting")
	}

	client := &http.Client{}
	requestURL := "https://api.spotify.com/v1/me/tracks?limit=1"
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.Join(err, errors.New("error retrieving favourite count"))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		retryAfter := resp.Header.Get("Retry-After")
		r, err := strconv.Atoi(retryAfter)
		if err == nil {
			fmt.Printf("RetrieveFavouriteCount rate limited, retry after: %s seconds\n", retryAfter)
			CacheStore.Set(LOCK_KEY, "true", time.Duration(r)*time.Second)
		}
		return 0, errors.New("rate limited by Spotify API")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("Spotify API returned status %d when retrieving favourite count", resp.StatusCode)
	}

	countObj := &SavedTracks{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, countObj)

	if err != nil {
		return 0, errors.Join(err, errors.New("error unmarshalling favourite count"))
	}

	return countObj.Total, nil
}
