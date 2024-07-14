package main

import (
	"database/sql"
	"fmt"
	"log"
)

func initDb() *sql.DB {
	dbFileName := getEnv("DB_FILE", "./shufflify.db")

	fmt.Printf("db file: %s\n", dbFileName)

	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (id text not null primary key, token text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	return db
}
