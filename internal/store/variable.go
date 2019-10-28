package store

import (
	"time"

	"github.com/mattn/go-sqlite3"
)

type Variable struct {
	Name        string
	Value       string
	Environment string
	Type        string
	Generator   string
	Timeout     int64
	Last        time.Time
}

func (v *Variable) Save() error {
	return StoreVariables([]Variable{*v})
}
func (v *Variable) Delete() error {
	_, err := globalDB.Exec("DELETE FROM variables WHERE name=$1 AND environment=$2",
		v.Name, v.Environment)
	if err != nil {
		log.Errorf("%+v\n", err)
		return ErrorVariableNotFound
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
			`INSERT OR REPLACE INTO variables
			(name, value, environment, type, generator, timeout, last) VALUES
			(:name, :value, :environment, :type, :generator, :timeout, :last)`,
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
		log.Errorf("%+v\n", err)
	}
	return variables
}
func GetVariablesByName(name string) []Variable {
	variables := []Variable{}
	if err := globalDB.Select(&variables,
		"SELECT * FROM variables WHERE name=$1", name); err != nil {
		log.Errorf("%+v\n", err)
	}
	return variables
}
func GetVariablesByEnvironment(environment string) []Variable {
	variables := []Variable{}
	if err := globalDB.Select(&variables,
		"SELECT * FROM variables WHERE environment=$1", environment); err != nil {
		log.Errorf("%+v\n", err)
	}
	return variables
}
func GetVariablesByType(typ string) []Variable {
	variables := []Variable{}
	if err := globalDB.Select(&variables,
		"SELECT * FROM variables WHERE type=$1", typ); err != nil {
		log.Errorf("%+v\n", err)
	}
	return variables
}
func GetVariableByNameAndEnvironment(name, environment string) (Variable, error) {
	variable := Variable{}
	if err := globalDB.Get(&variable,
		"SELECT * FROM variables WHERE name=$1 AND environment=$2",
		name, environment); err != nil {
		log.Errorf("%+v\n", err)
		return Variable{}, ErrorVariableNotFound
	}
	return variable, nil
}
func GetVariablesByEnvironmentAndType(environment, typ string) []Variable {
	variables := []Variable{}
	if err := globalDB.Select(&variables,
		"SELECT * FROM variables WHERE environment=$1 AND type=$2", environment, typ); err != nil {
		log.Errorf("%+v\n", err)
	}
	return variables
}
func GetVariablesByNameAndType(name, typ string) []Variable {
	variables := []Variable{}
	if err := globalDB.Select(&variables,
		"SELECT * FROM variables WHERE name=$1 AND type=$2", name, typ); err != nil {
		log.Errorf("%+v\n", err)
	}
	return variables
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
		timeout INT,
		last DATETIME,
		PRIMARY KEY (name, environment),
		FOREIGN KEY(environment) REFERENCES environments(name)
	);
	`

	_, err := globalDB.Exec(query)
	if err != nil {
		panic(err)
	}
}
