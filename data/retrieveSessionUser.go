package data

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/syeo66/go-shufflify/lib"
	. "github.com/syeo66/go-shufflify/types"
)

func RetrieveSessionUser(r *http.Request) (*User, error) {
	session, _ := lib.Store.Get(r, "user-session")

	userData := session.Values["user"]

	if userData == nil {
		return nil, errors.New("user not found")
	}

	user := &User{}
	err := json.Unmarshal(userData.([]byte), user)
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving user"))
	}

	return user, nil
}
