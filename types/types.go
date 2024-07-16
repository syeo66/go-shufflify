package types

import "time"

type Image struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type User struct {
	Id          string  `json:"id"`
	DisplayName string  `json:"display_name"`
	Images      []Image `json:"images"`
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type Page struct {
	Configuration *Configuration
	Player        *Player
	Queue         Queue
	User          User
}

type Artist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	Id      string   `json:"id"`
	Images  []Image  `json:"images"`
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

type Track struct {
	Id         string   `json:"id"`
	Uri        string   `json:"uri"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	Explicit   bool     `json:"explicit"`
	Name       string   `json:"name"`
	Popularity int      `json:"popularity"`
}

type Queue struct {
	CurrentlyPlaying Track   `json:"currently_playing"`
	Queue            []Track `json:"queue"`
}

type Device struct {
	Id       string `json:"id"`
	IsActive bool   `json:"isActive"`
	Name     string `json:"name"`
}

type Player struct {
	Device       *Device `json:"device"`
	RepeatState  string  `json:"repeat_state"`
	ShuffleState bool    `json:"shuffle_state"`
	ProgressMs   int     `json:"progress_ms"`
	IsPlaying    bool    `json:"is_playing"`
	Item         Track   `json:"item"`
}

type Configuration struct {
	UID         string     `json:"uid"`
	ActiveUntil *time.Time `json:"activeUntil"`
	IsActive    bool       `json:"isActive"`
}

type SavedTrack struct {
	AddedAt time.Time `json:"added_at"`
	Track   Track     `json:"track"`
}

type SavedTracks struct {
	Total int          `json:"total"`
	Items []SavedTrack `json:"items"`
}

type PlaylistTracks struct {
	Href  string `json:"href"`
	Total int    `json:"total"`
}

type Playlist struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Tracks      PlaylistTracks `json:"tracks"`
}

type Playlists struct {
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Total    int        `json:"total"`
	Items    []Playlist `json:"items"`
}
