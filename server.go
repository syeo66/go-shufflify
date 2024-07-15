package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/syeo66/go-shufflify/lib"
	"github.com/syeo66/go-shufflify/routes"
)

// TOOD
// - implement a token refresh mechanism
// - introduce queue manager (go routine)
// - prevent queue manager and frontend from fetching the same data twice

func main() {
	port := lib.GetEnv("PORT", "3333")

	db := lib.InitDb()
	defer db.Close()

	tmpl := make(map[string]*template.Template)
	tmpl["index.html"] = template.Must(template.ParseFiles("templates/base.html", "templates/index.html", "templates/player.html", "templates/queue.html"))
	tmpl["login.html"] = template.Must(template.ParseFiles("templates/base.html", "templates/login.html"))
	tmpl["player.html"] = template.Must(template.ParseFiles("templates/player.html"))
	tmpl["queue.html"] = template.Must(template.ParseFiles("templates/queue.html"))

	cssfs := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssfs))

	jsfs := http.FileServer(http.Dir("./js"))
	http.Handle("/js/", http.StripPrefix("/js/", jsfs))

	http.HandleFunc("/", routes.GetRoot(tmpl, db))
	http.HandleFunc("/callback", routes.GetCallback(db))
	http.HandleFunc("/login", routes.GetLogin(tmpl))
	http.HandleFunc("/logout", routes.GetLogout)
	http.HandleFunc("/player", routes.GetPlayer(tmpl, db))
	http.HandleFunc("/queue", routes.GetQueue(tmpl, db))

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
