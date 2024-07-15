package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func retrieveQueue(token string) (*Queue, error) {
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

	return cleanQueue, nil
}
