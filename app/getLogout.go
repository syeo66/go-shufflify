package main

import (
	"fmt"
	"net/http"
)

func getLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)

	session, _ := store.Get(r, "user-session")
	session.Values["user"] = ""
	_ = session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusFound)
}
