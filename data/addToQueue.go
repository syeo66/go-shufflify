package data

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/syeo66/go-shufflify/types"
)

func AddToQueue(token string, t Track) error {
	client := &http.Client{}
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=%s", t.Uri)
	req, _ := http.NewRequest(http.MethodPost, requestURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return errors.Join(err, errors.New("error while adding to queue"))
	}
	defer resp.Body.Close()

	return nil
}
