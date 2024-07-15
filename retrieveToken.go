package main

import "database/sql"

func retrieveToken(user *User, db *sql.DB) string {
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
