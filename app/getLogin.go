package main

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
)

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
