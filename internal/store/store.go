package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var globalDB *sql.DB

func InitDB() {
	db, err := sql.Open("sqlite3", "test.db?_foreign_keys=on")
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}

	// create targets table if not exists
	targetsTable := `
	CREATE TABLE IF NOT EXISTS targets(
		alias TEXT NOT NULL PRIMARY KEY,
		url TEXT
	);
	`
	_, err = db.Exec(targetsTable)
	if err != nil {
		panic(err)
	}

	// create requests table if not exists
	requestsTable := `
	CREATE TABLE IF NOT EXISTS requests(
		name TEXT NOT NULL PRIMARY KEY,
		target TEXT,
		method TEXT,
		FOREIGN KEY(target) REFERENCES targets(alias)
	);
	`

	_, err = db.Exec(requestsTable)
	if err != nil {
		panic(err)
	}

	globalDB = db
}
