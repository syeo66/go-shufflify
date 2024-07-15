package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	d "github.com/syeo66/go-shufflify/data"
	. "github.com/syeo66/go-shufflify/types"
)

func GetRoot(
	tmpl map[string]*template.Template,
	db *sql.DB,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		page := Page{}

		user, err := d.RetrieveSessionUser(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		page.User = *user
		token := d.RetrieveToken(user.Id, db)

		queue, err := d.RetrieveQueue(token)
		if err == nil && queue != nil {
			page.Queue = *queue
		} else {
			fmt.Println(err)
		}

		player, _ := d.RetrievePlayer(token)
		page.Player = player

		err = tmpl["index.html"].ExecuteTemplate(w, "base.html", page)
		if err != nil {
			fmt.Printf("error executing template: %s\n", err)
		}
	}
}
