package data

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	. "github.com/syeo66/go-shufflify/types"
)

func RefreshToken(uid string, refreshToken string, db *sql.DB) (string, error) {
	if refreshToken == "" {
		return "", errors.New("no refresh token")
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	auth := clientID + ":" + clientSecret
	auth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))

	data := url.Values{}
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	client := &http.Client{}
	requestURL := "https://accounts.spotify.com/api/token"
	req, _ := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(data.Encode()))
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Join(err, errors.New("error retrieving token"))
	}
	defer resp.Body.Close()

	bodyData := &AccessToken{}
	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, bodyData)

		if err != nil {
			return "", errors.Join(err, errors.New("error retrieving token"))
		}

		sqlStmt := `UPDATE users SET token = ?, expiry = ? WHERE id = ?`
		_, err = db.Exec(
			sqlStmt,
			bodyData.AccessToken,
			time.Now().Add(time.Second*time.Duration(bodyData.ExpiresIn)),
			uid,
		)

		if err != nil {
			return "", errors.Join(err, errors.New("error writing token"))
		}
	} else {
		return "", errors.New("error retrieving token")
	}

	return bodyData.AccessToken, nil
}
