package store

import (
	"fmt"
	"strings"
)

type EnvironmentType struct {
	Name string
}

func StoreEnvironments(envs []EnvironmentType) error {
	if len(envs) == 0 {
		return nil
	}
	request := `INSERT INTO environments(name) VALUES`
	var requestValues []string
	var values []interface{}
	for _, env := range envs {
		requestValues = append(requestValues, "(?)")
		values = append(values, env.Name)
	}
	request = fmt.Sprintf("%s %s;", request, strings.Join(requestValues, ","))

	stmt, err := globalDB.Prepare(request)
	if err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return err
	}
	return nil
}

func StoreEnvironment(env EnvironmentType) error {
	return StoreEnvironments([]EnvironmentType{env})
}

func GetAllEnvironments() []EnvironmentType {
	request := `SELECT name FROM environments`
	rows, err := globalDB.Query(request)
	if err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return []EnvironmentType{}
	}
	defer rows.Close()

	var result []EnvironmentType
	for rows.Next() {
		item := EnvironmentType{}
		err := rows.Scan(&item.Name)
		if err != nil {
			// TODO: log error
			fmt.Printf("error: %+v\n", err)
			return []EnvironmentType{}
		}
		result = append(result, item)
	}
	return result
}

func init() {
	if globalDB == nil {
		initDB()
	}
	// create requests table if not exists
	request := `
	CREATE TABLE IF NOT EXISTS environments(
		name TEXT NOT NULL PRIMARY KEY
	);
	`

	_, err := globalDB.Exec(request)
	if err != nil {
		panic(err)
	}
}
