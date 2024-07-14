package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)

	session, _ := store.Get(r, "user-session")

	requestURL := "https://accounts.spotify.com/api/token"
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
	req, _ := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(data.Encode()))
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)

	if resp.StatusCode == 200 {
		// TODO store access token in db
		// TODO initialize a session for the user
		fmt.Printf("redirecting to /\n")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fmt.Printf("error: %s\n", resp.Status)
	http.Redirect(w, r, "/login", http.StatusFound)
}
