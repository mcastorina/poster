package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var globalDB *sql.DB

func initDB() {
	db, err := sql.Open("sqlite3", "test.db?_foreign_keys=on")
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	globalDB = db
}
