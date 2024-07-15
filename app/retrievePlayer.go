package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func retrievePlayer(token string) (*Player, error) {
	client := &http.Client{}
	requestURL := "https://api.spotify.com/v1/me/player"
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving user"))
	}
	defer resp.Body.Close()

	player := &Player{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, player)

	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving user"))
	}

	return player, nil
}
