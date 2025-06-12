package data

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/syeo66/go-shufflify/types"
)

func AddToQueue(token string, t Track) error {
	_, found := CacheStore.Get(LOCK_KEY)
	if found {
		return errors.New("request locked due to rate limiting")
	}

	client := &http.Client{}
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=%s", t.Uri)
	req, _ := http.NewRequest(http.MethodPost, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return errors.Join(err, errors.New("error while adding to queue"))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		retryAfter := resp.Header.Get("Retry-After")
		r, err := strconv.Atoi(retryAfter)
		if err == nil {
			fmt.Printf("AddToQueue rate limited, retry after: %s seconds\n", retryAfter)
			CacheStore.Set(LOCK_KEY, "true", time.Duration(r)*time.Second)
		}
		return errors.New("rate limited by Spotify API")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Spotify API returned status %d when adding to queue", resp.StatusCode)
	}

	return nil
}
