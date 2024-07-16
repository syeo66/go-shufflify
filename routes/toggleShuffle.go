package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	d "github.com/syeo66/go-shufflify/data"
)

func ToggleShuffle(tmpl map[string]*template.Template, db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)

		if r.Method != "GET" && r.Method != "POST" {
			http.NotFound(w, r)
			return
		}

		user, _ := d.RetrieveSessionUser(r)

		if r.Method == "POST" {
			sqlStmt := `
      UPDATE users 
      SET IsActive = CASE IsActive
                     WHEN 0 THEN 1 
                     ELSE 0 
                     END 
      WHERE id = ?
      `
			_, err := db.Exec(
				sqlStmt,
				user.Id,
			)
			if err != nil {
				fmt.Println(err)
			}
		}

		configuration, err := d.RetrieveConfig(user.Id, db)
		if err != nil {
			fmt.Println(err)
		}

		err = tmpl["tools.html"].ExecuteTemplate(w, "tools.html", configuration)
		if err != nil {
			fmt.Println(err)
		}
	}
}
