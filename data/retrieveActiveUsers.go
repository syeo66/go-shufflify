package data

import (
	"database/sql"
	"errors"
	"fmt"
)

func RetrieveActiveUsers(db *sql.DB) ([]string, error) {
	rows, err := db.Query("select id from users where isActive = true")
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving active users"))
	}
	defer rows.Close()

	resp := []string{}

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		resp = append(resp, id)
	}

	err = rows.Err()
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving active users"))
	}

	return resp, nil
}
