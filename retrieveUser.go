package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func retrieveUser(token string) (*User, error) {
	client := &http.Client{}
	requestURL := "https://api.spotify.com/v1/me"
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving user"))
	}
	defer resp.Body.Close()

	user := &User{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, user)

	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving user"))
	}

	return user, nil
}
