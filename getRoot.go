package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func getRoot(
	tmpl map[string]*template.Template,
	db *sql.DB,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)

		session, _ := store.Get(r, "user-session")

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		page := Page{}

		userData := session.Values["user"]

		if userData == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user := &User{}
		err := json.Unmarshal(userData.([]byte), user)
		if err == nil {
			page.User = *user
		}

		token := retrieveToken(user, db)

		queue, err := retrieveQueue(token)
		if err == nil && queue != nil {
			page.Queue = *queue
		} else {
			fmt.Println(err)
		}

		player, err := retrievePlayer(token)
		if err == nil && player != nil {
			page.Player = *player
		} else {
			fmt.Println(err)
		}

		err = tmpl["index.html"].ExecuteTemplate(w, "base.html", page)
		if err != nil {
			fmt.Printf("error executing template: %s\n", err)
		}
	}
}
