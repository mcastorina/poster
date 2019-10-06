package store

import (
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type Variable struct {
	Name        string
	Value       string
	Environment string
	Type        string
	Generator   string
}

func (v *Variable) Save() error {
	return StoreVariables([]Variable{*v})
}
func (v *Variable) Delete() error {
	_, err := globalDB.Exec("DELETE FROM variables WHERE name=$1 AND environment=$2",
		v.Name, v.Environment)
	if err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return ErrorVariableNotFound
	}
	return nil
}
func (v *Variable) Update() error {
	_, err := globalDB.NamedExec(
		`UPDATE variables SET value=:value, type=:type, generator=:generator
		WHERE name=:name AND environment=:environment`, v)
	if err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return err
	}
	return nil
}

func StoreVariables(variables []Variable) error {
	if len(variables) == 0 {
		return nil
	}
	tx := globalDB.MustBegin()

	for _, variable := range variables {
		if _, err := tx.NamedExec(
			`INSERT INTO variables (name, value, environment, type, generator)
			VALUES (:name, :value, :environment, :type, :generator)`,
			&variable); err != nil {

			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.Code == sqlite3.ErrConstraint {
					if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
						return ErrorVariableExists
					}
					return ErrorEnvironmentNotFound
				}
				return ErrorUnknown
			}
			// Should not reach
			return err
		}
	}

	return tx.Commit()
}

func GetAllVariables() []Variable {
	variables := []Variable{}
	if err := globalDB.Select(&variables, "SELECT * FROM variables"); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
	}
	return variables
}
func GetVariablesByEnvironment(environment string) []Variable {
	variables := []Variable{}
	if err := globalDB.Select(&variables,
		"SELECT * FROM variables WHERE environment=$1", environment); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
	}
	return variables
}
func GetVariablesByName(name string) []Variable {
	variables := []Variable{}
	if err := globalDB.Select(&variables,
		"SELECT * FROM variables WHERE name=$1", name); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
	}
	return variables
}
func GetVariableByNameAndEnvironment(name, environment string) (Variable, error) {
	variable := Variable{}
	if err := globalDB.Get(&variable,
		"SELECT * FROM variables WHERE name=$1 AND environment=$2",
		name, environment); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return Variable{}, ErrorVariableNotFound
	}
	return variable, nil
}

func init() {
	if globalDB == nil {
		initDB()
	}
	// create requests table if not exists
	query := `
	CREATE TABLE IF NOT EXISTS variables(
		name TEXT NOT NULL,
		value TEXT,
		environment TEXT NOT NULL,
		type TEXT NOT NULL,
		generator TEXT,
		PRIMARY KEY (name, environment),
		FOREIGN KEY(environment) REFERENCES environments(name)
	);
	`

	_, err := globalDB.Exec(query)
	if err != nil {
		panic(err)
	}
}
