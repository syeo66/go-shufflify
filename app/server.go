package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

// TOOD
// - display and update queue
// - introduce queue manager (go routine)
// - prevent queue manager and frontend from fetching the same data twice

func main() {
	port := getEnv("PORT", "3333")

	db := initDb()
	defer db.Close()

	tmpl := make(map[string]*template.Template)
	tmpl["index.html"] = template.Must(template.ParseFiles("templates/base.html", "templates/index.html", "templates/player.html"))
	tmpl["login.html"] = template.Must(template.ParseFiles("templates/base.html", "templates/login.html"))
	tmpl["player.html"] = template.Must(template.ParseFiles("templates/player.html"))

	cssfs := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssfs))

	jsfs := http.FileServer(http.Dir("./js"))
	http.Handle("/js/", http.StripPrefix("/js/", jsfs))

	http.HandleFunc("/", getRoot(tmpl, db))
	http.HandleFunc("/callback", getCallback(db))
	http.HandleFunc("/login", getLogin(tmpl))
	http.HandleFunc("/logout", getLogout)
	http.HandleFunc("/player", getPlayer(tmpl, db))

	fmt.Printf("starting server port %s\n", port)
	fmt.Printf("open http://localhost:%s/\n", port)

	err := http.ListenAndServe(":"+port, nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
