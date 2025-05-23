package routes

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	d "github.com/syeo66/go-shufflify/data"
	"github.com/syeo66/go-shufflify/lib"
	types "github.com/syeo66/go-shufflify/types"
)

func GetCallback(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)

		session, _ := lib.Store.Get(r, "user-session")

		redirectURI := fmt.Sprintf("http://%s/callback", r.Host)

		clientID := os.Getenv("SPOTIFY_CLIENT_ID")
		clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if session.Values["state"] != state {
			session.Values["state"] = ""
			err := session.Save(r, w)
			if err != nil {
				fmt.Printf("error saving session: %s\n", err)
				return
			}

			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		auth := clientID + ":" + clientSecret
		auth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))

		data := url.Values{}
		data.Set("code", code)
		data.Set("redirect_uri", redirectURI)
		data.Set("grant_type", "authorization_code")

		client := &http.Client{}
		requestURL := "https://accounts.spotify.com/api/token"
		req, _ := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(data.Encode()))
		req.Header.Add("Authorization", auth)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			body, _ := io.ReadAll(resp.Body)
			bodyData := &types.AccessToken{}
			err = json.Unmarshal(body, bodyData)

			if err != nil {
				fmt.Printf("error: %s\n", err)
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			user, _ := d.RetrieveUser(bodyData.AccessToken)
			config, _ := d.RetrieveConfig(user.Id, db)

			disallowNewUsers := os.Getenv("DISALLOW_NEW_USERS")
			if config == nil && disallowNewUsers == "true" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			if config == nil {
				sqlStmt := `
      REPLACE INTO users (id, token, refreshToken, expiry, isActive) VALUES (?, ?, ?, ?, ?);
      `
				_, err = db.Exec(
					sqlStmt,
					user.Id,
					bodyData.AccessToken,
					bodyData.RefreshToken,
					time.Now().Add(time.Second*time.Duration(bodyData.ExpiresIn)),
					0,
				)

				if err != nil {
					fmt.Printf("error: %s\n", err)
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
			}

			jsonUser, _ := json.Marshal(user)
			session.Values["user"] = jsonUser
			err = session.Save(r, w)
			if err != nil {
				fmt.Printf("error saving session: %s\n", err)
			}

			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		fmt.Printf("error: %s\n", resp.Status)
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
