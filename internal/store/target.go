package store

import (
	"fmt"
)

type Target struct {
	Alias string
	URL   string
}

func StoreTargets(targets []Target) error {
	if len(targets) == 0 {
		return nil
	}
	tx := globalDB.MustBegin()

	for _, target := range targets {
		tx.NamedExec("INSERT INTO targets (alias, url) VALUES (:alias, :url)", &target)
	}

	return tx.Commit()
}

func StoreTarget(target Target) error {
	return StoreTargets([]Target{target})
}

func GetAllTargets() []Target {
	targets := []Target{}
	if err := globalDB.Select(&targets, "SELECT * FROM targets"); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
	}
	return targets
}

func GetTargetByAlias(alias string) (Target, error) {
	target := Target{}
	if err := globalDB.Get(&target, "SELECT * FROM targets WHERE alias=$1", alias); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return Target{}, err
	}
	return target, nil
}

func GetTargetByURL(url string) (Target, error) {
	target := Target{}
	if err := globalDB.Get(&target, "SELECT * FROM targets WHERE url=$1", url); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return Target{}, err
	}
	return target, nil
}

func init() {
	if globalDB == nil {
		initDB()
	}
	// create targets table if not exists
	request := `
	CREATE TABLE IF NOT EXISTS targets(
		alias TEXT NOT NULL PRIMARY KEY,
		url TEXT
	);
	`
	_, err := globalDB.Exec(request)
	if err != nil {
		panic(err)
	}
}
