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

func RetrieveQueue(token string) (*Queue, error) {
	key := "RetrieveQueue" + token

	_, found := CacheStore.Get(LOCK_KEY)

	if found {
		return nil, errors.New("request locked")
	}

	cachedQueue, found := CacheStore.Get(key)
	if found {
		return cachedQueue.(*Queue), nil
	}

	client := &http.Client{}
	requestURL := "https://api.spotify.com/v1/me/player/queue"
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving queue"))
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

	queue := &Queue{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, queue)

	cleanQueue := &Queue{
		CurrentlyPlaying: queue.CurrentlyPlaying,
		Queue:            []Track{},
	}

	seen := map[string]bool{}

	// For some reason the spotify api returns the same track multiple times
	// even if it is not in the queue only once. This makes the API to
	// retreive the queue mostly useless.

	// TODO remove this later when Spotify works properly (which might be never)
	for _, track := range queue.Queue {
		if seen[track.Id] {
			continue
		}
		cleanQueue.Queue = append(cleanQueue.Queue, track)
		seen[track.Id] = true
	}

	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving queue"))
	}

	CacheStore.Set(key, cleanQueue, 5*time.Second)

	return cleanQueue, nil
}
