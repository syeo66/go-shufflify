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

const LOCK_KEY = "request-locked"

func RetrievePlayer(token string) (*Player, error) {
	key := "RetrievePlayer" + token

	_, found := CacheStore.Get(LOCK_KEY)

	if found {
		return nil, errors.New("request locked")
	}

	cachedPlayer, found := CacheStore.Get(key)
	if found {
		return cachedPlayer.(*Player), nil
	}

	client := &http.Client{}
	requestURL := "https://api.spotify.com/v1/me/player"
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving user"))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 {
		retryAfter := resp.Header.Get("Retry-After")
		r, err := strconv.Atoi(retryAfter)
		if err == nil {
			fmt.Printf("Retry-After: %s\n", retryAfter)
			CacheStore.Set(LOCK_KEY, "true", time.Duration(r)*time.Second)
		}
		return nil, errors.New("Throttle")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Spotify API returned status %d when retrieving player", resp.StatusCode)
	}

	player := &Player{}

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, player)

	if err != nil {

		return nil, errors.Join(err, errors.New("error unmarshalling user"))
	}

	CacheStore.Set(key, player, 2*time.Second)

	return player, nil
}
