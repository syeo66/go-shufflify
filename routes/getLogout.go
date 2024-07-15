package routes

import (
	"fmt"
	"net/http"

	"github.com/syeo66/go-shufflify/lib"
)

func GetLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)

	session, _ := lib.Store.Get(r, "user-session")
	session.Values["user"] = nil
	_ = session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusFound)
}
