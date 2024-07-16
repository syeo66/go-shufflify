package data

import (
	"database/sql"
	"errors"
	"time"

	. "github.com/syeo66/go-shufflify/types"
)

func RetrieveConfig(uid string, db *sql.DB) (*Configuration, error) {
	stmt, err := db.Prepare("select isActive, activeUntil from users where id = ?")
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving config"))
	}
	defer stmt.Close()

	var isActive bool
	var activeUntil *time.Time
	err = stmt.QueryRow(uid).Scan(&isActive, &activeUntil)
	if err != nil {
		return nil, errors.Join(err, errors.New("error retrieving config"))
	}

	return &Configuration{
		UID:         uid,
		IsActive:    isActive,
		ActiveUntil: activeUntil,
	}, nil
}
