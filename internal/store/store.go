package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var globalDB *sqlx.DB

func initDB() {
	db, err := sqlx.Open("sqlite3", "test.db?_foreign_keys=on")
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	globalDB = db
}
