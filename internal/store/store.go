package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var globalDB *sql.DB

type Target struct {
	Alias string
	URL   string
}

func InitDB() {
	db, err := sql.Open("sqlite3", "test.db")
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

	globalDB = db
}

func StoreTarget(target Target) error {
	request := `
	INSERT OR REPLACE INTO targets(
		alias,
		url
	) values(?, ?)
	`

	stmt, err := globalDB.Prepare(request)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(target.Alias, target.URL)
	if err != nil {
		return err
	}
	return nil
}

func GetAllTargets() []Target {
	request := `
	SELECT alias,url FROM targets;
	`
	rows, err := globalDB.Query(request)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []Target
	for rows.Next() {
		item := Target{}
		err := rows.Scan(&item.Alias, &item.URL)
		if err != nil {
			panic(err)
		}
		result = append(result, item)
	}
	return result
}
