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
	users, err := d.RetrieveActiveUsers(db)
	if err != nil {
		fmt.Printf("Error retrieving active users: %v\n", err)
		return
	}

	for _, uid := range users {
		if err := processUserQueue(uid, db); err != nil {
			fmt.Printf("Error processing queue for user %s: %v\n", uid, err)
		}
	}
}

func processUserQueue(uid string, db *sql.DB) error {
	token, err := d.RetrieveToken(uid, db)
	if err != nil {
		return fmt.Errorf("failed to retrieve token for user %s: %w", uid, err)
	}
	if token == "" {
		return fmt.Errorf("received empty token for user %s", uid)
	}

	player, err := d.RetrievePlayer(token)
	if err != nil {
		return fmt.Errorf("failed to retrieve player for user %s: %w", uid, err)
	}
	if player == nil {
		return fmt.Errorf("no active player found for user %s", uid)
	}

	cacheKey := "playedTracks" + uid + player.Item.Id
	d.CacheStore.Set(cacheKey, true, 12*time.Hour)

	queue, err := d.RetrieveQueue(token)
	if err != nil {
		return fmt.Errorf("failed to retrieve queue for user %s: %w", uid, err)
	}
	if queue == nil || queue.Queue == nil {
		return fmt.Errorf("no queue found for user %s", uid)
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
		return nil // Queue has enough songs, nothing to do
	}

	favCount, err := d.RetrieveFavouriteCount(token, db)
	if err != nil {
		return fmt.Errorf("failed to retrieve favourite count for user %s: %w", uid, err)
	}

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
		return fmt.Errorf("user %s has no tracks available (favourites: %d, playlists: %d)", uid, favCount, playlistsCount)
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(totalCount)))
	if err != nil {
		return fmt.Errorf("failed to generate random number for user %s: %w", uid, err)
	}
	num := int(n.Int64())

	var track types.Track

	if num < favCount {
		fmt.Println("Fav")
		t, err := d.RetrieveNthSongFromFavourites(token, num)
		if err != nil {
			return fmt.Errorf("failed to retrieve favourite song %d for user %s: %w", num, uid, err)
		}
		if t == nil {
			return fmt.Errorf("no favourite song found at position %d for user %s", num, uid)
		}
		track = *t
	} else {
		num = num - favCount

		for _, p := range playlists {
			if num < p.Tracks.Total {
				fmt.Println("playlist")
				t, err := d.RetrieveNthSongFromPlaylist(token, p, num)
			if err != nil {
				return fmt.Errorf("failed to retrieve song %d from playlist %s for user %s: %w", num, p.Name, uid, err)
			}
			if t == nil {
				return fmt.Errorf("no song found at position %d in playlist %s for user %s", num, p.Name, uid)
			}
				track = *t
				break
			} else {
				num = num - p.Tracks.Total
			}
		}
	}

	err = d.AddToQueue(token, track)
	if err != nil {
		return fmt.Errorf("failed to add track %s to queue for user %s: %w", track.Name, uid, err)
	}

	return nil
}
