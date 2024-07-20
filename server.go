package main

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	d "github.com/syeo66/go-shufflify/data"
	"github.com/syeo66/go-shufflify/lib"
	"github.com/syeo66/go-shufflify/routes"
	. "github.com/syeo66/go-shufflify/types"
)

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

		if token == "" {
			continue
		}

		player, _ := d.RetrievePlayer(token)

		if player == nil {
			continue
		}

		cacheKey := "playedTracks" + token + player.Item.Id
		d.CacheStore.Set(cacheKey, true, 12*time.Hour)

		queue, _ := d.RetrieveQueue(token)

		if queue == nil || queue.Queue == nil {
			continue
		}

		queueList := []Track{}

		// remove played tracks from queue
		for _, t := range queue.Queue {
			cacheKey := "playedTracks" + token + t.Id
			_, found := d.CacheStore.Get(cacheKey)

			if !found {
				queueList = append(queueList, t)
			}
		}

		if len(queueList) > 3 {
			continue
		}

		favCount := d.RetrieveFavouriteCount(token, db)

		playlists, err := d.RetrievePlaylists(token, db)
		if err != nil {
			fmt.Println(err)
		}

		playlistsCount := 0
		for _, p := range playlists {
			playlistsCount += p.Tracks.Total
		}

		totalCount := favCount + playlistsCount

		if totalCount <= 0 {
			continue
		}

		n, _ := rand.Int(rand.Reader, big.NewInt(int64(totalCount)))
		num := int(n.Int64())
		var track Track

		if num < favCount {
			t, _ := d.RetrieveNthSongFromFavourites(token, num)
			if t == nil {
				continue
			}
			track = *t
		} else {
			num = num - favCount
			for _, p := range playlists {
				if num < p.Tracks.Total {
					t, _ := d.RetrieveNthSongFromPlaylist(token, p, num)
					track = *t
					break
				} else {
					num = num - p.Tracks.Total
				}
			}
		}

		err = d.AddToQueue(token, track)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
