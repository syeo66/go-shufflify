package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	. "github.com/syeo66/go-shufflify/types"
)

func RetrieveQueue(token string) (*Queue, error) {
	key := "RetrieveQueue" + token

	cachedQueue, found := cacheStore.Get(key)
	if found {
		return cachedQueue.(*Queue), nil
	}

	client := &http.Client{}
	requestURL := "https://api.spotify.com/v1/me/player/queue"
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving user"))
	}
	defer resp.Body.Close()

	queue := &Queue{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, queue)

	cleanQueue := &Queue{
		CurrentlyPlaying: queue.CurrentlyPlaying,
		Queue:            []Track{},
	}

	// for some reason the spotify api returns the same track multiple times
	// even if it is not in the queue only once.
	// TODO remove this later when Spotify works properly (which might be never)
	var prevId string
	for _, track := range queue.Queue {
		if prevId == track.Id {
			continue
		}
		cleanQueue.Queue = append(cleanQueue.Queue, track)
		prevId = track.Id
	}

	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving user"))
	}

	cacheStore.Set(key, cleanQueue, 5*time.Second)

	return cleanQueue, nil
}
