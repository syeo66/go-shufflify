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

		token, err := d.RetrieveToken(user.Id, db)
		if err != nil {
			fmt.Printf("Error retrieving token: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		player, err := d.RetrievePlayer(token)
		if err != nil {
			fmt.Printf("Error retrieving player: %v\n", err)
		}

		err = tmpl["player.html"].ExecuteTemplate(w, "player.html", player)
		if err != nil {
			fmt.Println(err)
		}
	}
}
