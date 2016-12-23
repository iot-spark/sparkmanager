package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)

	if err != nil {
		log.Fatalf("Init sqlite database error: %#v\n", err)
	}

	if db == nil {
		log.Fatalf("Could not initialize sqlite database: %#v\n", filepath)
	}

	return db
}
