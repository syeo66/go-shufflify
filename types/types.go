package types

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
	AccessToken string `json:"access_token"`
	Tcope       string `json:"scope"`
}

type Page struct {
	User   User
	Queue  Queue
	Player *Player
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
