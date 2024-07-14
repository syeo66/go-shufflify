package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func getRoot(tmpl map[string]*template.Template) func(http.ResponseWriter, *http.Request) {
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

		err = tmpl["index.html"].ExecuteTemplate(w, "base.html", page)
		if err != nil {
			fmt.Printf("error executing template: %s\n", err)
		}
	}
}
