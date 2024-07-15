package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	d "github.com/syeo66/go-shufflify/data"
)

func GetPlayer(tmpl map[string]*template.Template, db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)

		user, err := d.RetrieveSessionUser(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		token := d.RetrieveToken(user.Id, db)
		player, _ := d.RetrievePlayer(token)

		err = tmpl["player.html"].ExecuteTemplate(w, "player.html", player)
		if err != nil {
			fmt.Println(err)
		}
	}
}
