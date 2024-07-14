package main

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
