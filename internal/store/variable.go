package store

import (
	"fmt"
)

type Variable struct {
	Name        string
	Value       string
	Environment string
	Type        string
	Generator   string
}

func StoreVariables(variables []Variable) error {
	if len(variables) == 0 {
		return nil
	}
	tx := globalDB.MustBegin()

	for _, variable := range variables {
		tx.NamedExec(
			`INSERT INTO variables (name, value, environment, type, generator)
			VALUES (:name, :value, :environment, :type, :generator)`,
			&variable)
	}

	return tx.Commit()
}
func StoreVariable(variable Variable) error {
	return StoreVariables([]Variable{variable})
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
func GetVariableByName(name string) (Variable, error) {
	variable := Variable{}
	if err := globalDB.Get(&variable,
		"SELECT * FROM variables WHERE name=$1", name); err != nil {
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
