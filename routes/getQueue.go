package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	d "github.com/syeo66/go-shufflify/data"
	"github.com/syeo66/go-shufflify/lib"
	. "github.com/syeo66/go-shufflify/types"
)

func GetQueue(tmpl map[string]*template.Template, db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)
		session, _ := lib.Store.Get(r, "user-session")

		userData := session.Values["user"]

		if userData == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user := &User{}
		err := json.Unmarshal(userData.([]byte), user)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		token := d.RetrieveToken(user, db)

		queue, err := d.RetrieveQueue(token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = tmpl["queue.html"].ExecuteTemplate(w, "queue.html", queue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
