package data

import (
	"database/sql"
	"fmt"
	"time"
)

func RetrieveToken(uid string, db *sql.DB) string {
	stmt, err := db.Prepare("select token, refreshToken, expiry from users where id = ?")
	if err != nil {
		return ""
	}
	defer stmt.Close()

	var token, refreshToken string
	var expiry time.Time
	err = stmt.QueryRow(uid).Scan(&token, &refreshToken, &expiry)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if expiry.Add(-5 * time.Minute).Before(time.Now()) {
		fmt.Println("try token refresh")

		token, err = RefreshToken(uid, refreshToken, db)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	return token
}
