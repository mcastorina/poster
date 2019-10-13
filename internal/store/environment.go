package store

import (
	"github.com/mattn/go-sqlite3"
)

type Environment struct {
	Name string
}

func (e *Environment) Save() error {
	return StoreEnvironments([]Environment{*e})
}
func (e *Environment) Delete() error {
	_, err := globalDB.Exec("DELETE FROM environments WHERE name=$1", e.Name)
	if err != nil {
		log.Errorf("%+v\n", err)
	}

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return ErrorEnvironmentInUse
		} else if sqliteErr.Code == sqlite3.ErrError {
			return ErrorEnvironmentNotFound
		} else {
			return ErrorUnknown
		}
	}
	return nil
}

func StoreEnvironments(envs []Environment) error {
	if len(envs) == 0 {
		return nil
	}
	tx := globalDB.MustBegin()

	for _, env := range envs {
		if _, err := tx.NamedExec(
			"INSERT OR REPLACE INTO environments (name) VALUES (:name)",
			&env); err != nil {

			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.Code == sqlite3.ErrConstraint {
					return ErrorEnvironmentExists
				}
				return ErrorUnknown
			}
		}
	}

	return tx.Commit()
}

func GetAllEnvironments() []Environment {
	envs := []Environment{}
	if err := globalDB.Select(&envs, "SELECT * FROM environments"); err != nil {
		log.Errorf("%+v\n", err)
	}
	return envs
}
func GetEnvironmentByName(name string) (Environment, error) {
	environment := Environment{}
	if err := globalDB.Get(&environment, "SELECT * FROM environments WHERE name=$1", name); err != nil {
		log.Errorf("%+v\n", err)
		return Environment{}, ErrorEnvironmentNotFound
	}
	return environment, nil
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
