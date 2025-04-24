package main

import (
	"crypto/rand"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
	d "github.com/syeo66/go-shufflify/data"
	"github.com/syeo66/go-shufflify/lib"
	"github.com/syeo66/go-shufflify/routes"
	types "github.com/syeo66/go-shufflify/types"
)

//go:embed css/*
var css embed.FS

//go:embed js/*
var js embed.FS

//go:embed images/*
var images embed.FS

func main() {
	port := lib.GetEnv("PORT", "3333")

	db := lib.InitDb()
	defer db.Close()

	worker(db)

	tmpl := lib.PrepareTemplates()

	cssfs := http.FileServer(http.FS(css))
	http.Handle("/css/", cssfs)

	jsfs := http.FileServer(http.FS(js))
	http.Handle("/js/", jsfs)

	imgfs := http.FileServer(http.FS(images))
	http.Handle("/images/", imgfs)

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
		processUserQueue(uid, db)
	}
}

func processUserQueue(uid string, db *sql.DB) {
	token := d.RetrieveToken(uid, db)

	if token == "" {
		return
	}

	player, _ := d.RetrievePlayer(token)

	if player == nil {
		return
	}

	cacheKey := "playedTracks" + uid + player.Item.Id
	d.CacheStore.Set(cacheKey, true, 12*time.Hour)

	queue, _ := d.RetrieveQueue(token)

	if queue == nil || queue.Queue == nil {
		return
	}

	queueList := []types.Track{}

	// remove played tracks from queue
	for _, t := range queue.Queue {
		cacheKey := "playedTracks" + uid + t.Id
		_, found := d.CacheStore.Get(cacheKey)

		if !found {
			queueList = append(queueList, t)
		}
	}

	if len(queueList) > 3 {
		return
	}

	favCount := d.RetrieveFavouriteCount(token, db)

	playlists, err := d.RetrievePlaylists(token, db)
	if err != nil {
		fmt.Println(err)
	}

	sort.Slice(playlists, func(i, j int) bool {
		return playlists[i].Tracks.Total > playlists[j].Tracks.Total
	})

	playlistsCount := 0
	for _, p := range playlists {
		playlistsCount += p.Tracks.Total
	}

	totalCount := favCount + playlistsCount

	if totalCount <= 0 {
		return
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(totalCount)))
	num := int(n.Int64())

	var track types.Track

	if num < favCount {
		fmt.Println("Fav")
		t, _ := d.RetrieveNthSongFromFavourites(token, num)
		if t == nil {
			return
		}
		track = *t
	} else {
		num = num - favCount

		for _, p := range playlists {
			if num < p.Tracks.Total {
				fmt.Println("playlist")
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
