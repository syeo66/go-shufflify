package main

import (
	"fmt"
	"net/http"
)

func getLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
	http.Redirect(w, r, "/login", http.StatusFound)
}
