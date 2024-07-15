package data

import (
	"database/sql"
)

func RetrieveToken(uid string, db *sql.DB) string {
	stmt, err := db.Prepare("select token from users where id = ?")
	if err != nil {
		return ""
	}
	defer stmt.Close()

	var token string
	err = stmt.QueryRow(uid).Scan(&token)
	if err != nil {
		return ""
	}

	return token
}
