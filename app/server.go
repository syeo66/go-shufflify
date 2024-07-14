package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

func main() {
	tmpl := make(map[string]*template.Template)
	tmpl["index.html"] = template.Must(template.ParseFiles("templates/base.html", "templates/index.html", "templates/player.html"))
	tmpl["login.html"] = template.Must(template.ParseFiles("templates/base.html", "templates/login.html"))

	http.HandleFunc("/", getRoot(tmpl))
	http.HandleFunc("/login", getLogin(tmpl))
	http.HandleFunc("/callback", getCallback)

	fmt.Printf("starting server port 3333\n")
	fmt.Printf("open http://localhost:3333/\n")
	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getCallback(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)

	requestURL := "https://accounts.spotify.com/api/token"
	redirectURI := fmt.Sprintf("http://%s/callback", r.Host)
	code := r.URL.Query().Get("code")
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	// TODO state check

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

func getLogin(tmpl map[string]*template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v+\n", r)
		fmt.Printf("%s /login\n", r.Method)

		state := generateRandomString(16)

		if r.Method == "POST" {
			spotifyUrl := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&state=%s&redirect_uri=http://%s/callback&scope=%s", os.Getenv("SPOTIFY_CLIENT_ID"), state, r.Host, "user-read-private")

			http.Redirect(w, r, spotifyUrl, http.StatusFound)
			return
		}

		err := tmpl["login.html"].ExecuteTemplate(w, "base.html", nil)
		if err != nil {
			fmt.Printf("error executing template: %s\n", err)
		}
	}
}

func getRoot(tmpl map[string]*template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		err := tmpl["index.html"].ExecuteTemplate(w, "base.html", nil)
		if err != nil {
			fmt.Printf("error executing template: %s\n", err)
		}
	}
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}
