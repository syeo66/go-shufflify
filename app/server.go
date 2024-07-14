package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"text/template"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

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
