package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/syeo66/go-shufflify/lib"
)

func GetLogin(tmpl map[string]*template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s /login\n", r.Method)

		session, _ := lib.Store.Get(r, "user-session")
		state := lib.GenerateRandomString(16)
		session.Values["state"] = state

		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		scopes := []string{
			"playlist-read-private",
			"user-library-read",
			"user-modify-playback-state",
			"user-read-currently-playing",
			"user-read-playback-state",
			"user-read-private",
		}
		scope := strings.Join(scopes, " ")
		if r.Method == "POST" {
			spotifyUrl := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&state=%s&redirect_uri=http://%s/callback&scope=%s", os.Getenv("SPOTIFY_CLIENT_ID"), state, r.Host, scope)

			http.Redirect(w, r, spotifyUrl, http.StatusFound)
			return
		}

		err = tmpl["login.html"].ExecuteTemplate(w, "base.html", nil)
		if err != nil {
			http.Error(w, "error executing template: login.html\n", http.StatusInternalServerError)
		}
	}
}
