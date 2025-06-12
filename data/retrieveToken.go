package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func RetrieveToken(uid string, db *sql.DB) (string, error) {
	stmt, err := db.Prepare("select token, refreshToken, expiry from users where id = ?")
	if err != nil {
		return "", errors.Join(err, errors.New("error preparing token query"))
	}
	defer stmt.Close()

	var token, refreshToken string
	var expiry time.Time
	err = stmt.QueryRow(uid).Scan(&token, &refreshToken, &expiry)
	if err != nil {
		return "", errors.Join(err, fmt.Errorf("error retrieving token for user %s", uid))
	}

	if expiry.Add(-5 * time.Minute).Before(time.Now()) {
		fmt.Println("try token refresh")

		token, err = RefreshToken(uid, refreshToken, db)
		if err != nil {
			return "", fmt.Errorf("error refreshing token for user %s: %w", uid, err)
		}
	}

	return token, nil
}
