package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	d "github.com/syeo66/go-shufflify/data"
	"github.com/syeo66/go-shufflify/lib"
	"github.com/syeo66/go-shufflify/routes"
)

// TOOD
// - introduce queue manager (go routine)
// - prevent queue manager and frontend from fetching the same data twice

func main() {
	port := lib.GetEnv("PORT", "3333")

	db := lib.InitDb()
	defer db.Close()

	worker(db)

	tmpl := lib.PrepareTemplates()

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
	http.HandleFunc("/toggle-shuffle", routes.ToggleShuffle(tmpl, db))

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

func worker(db *sql.DB) {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				queueManager(db)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func queueManager(db *sql.DB) {
	users, _ := d.RetrieveActiveUsers(db)

	for _, uid := range users {
		token := d.RetrieveToken(uid, db)
		queue, _ := d.RetrieveQueue(token)

		if queue.Queue == nil || len(queue.Queue) > 3 {
			continue
		}

		favCount := d.RetrieveFavouriteCount(token, db)
		fmt.Printf("favourites: %d\n", favCount)

		fmt.Println("add to queue")
	}
}
