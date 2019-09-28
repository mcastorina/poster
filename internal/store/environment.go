package store

import (
	"fmt"
)

type Environment struct {
	Name string
}

func StoreEnvironments(envs []Environment) error {
	if len(envs) == 0 {
		return nil
	}
	tx := globalDB.MustBegin()

	for _, env := range envs {
		tx.NamedExec("INSERT INTO environments (name) VALUES (:name)", &env)
	}

	return tx.Commit()
}

func StoreEnvironment(env Environment) error {
	return StoreEnvironments([]Environment{env})
}

func GetAllEnvironments() []Environment {
	envs := []Environment{}
	if err := globalDB.Select(&envs, "SELECT * FROM environments"); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
	}
	return envs
}

func init() {
	if globalDB == nil {
		initDB()
	}
	// create environments table if not exists
	query := `
	CREATE TABLE IF NOT EXISTS environments(
		name TEXT NOT NULL PRIMARY KEY
	);
	`

	_, err := globalDB.Exec(query)
	if err != nil {
		panic(err)
	}
}
