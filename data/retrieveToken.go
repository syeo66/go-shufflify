package data

import (
	"database/sql"
	. "github.com/syeo66/go-shufflify/types"
)

func RetrieveToken(user *User, db *sql.DB) string {
	id := user.Id

	stmt, err := db.Prepare("select token from users where id = ?")
	if err != nil {
		return ""
	}
	defer stmt.Close()

	var token string
	err = stmt.QueryRow(id).Scan(&token)
	if err != nil {
		return ""
	}

	return token
}
