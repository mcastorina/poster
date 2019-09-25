package store

import (
	"fmt"
	"strings"

	"github.com/mcastorina/poster/internal/models"
)

func StoreRequests(requests []models.Request) error {
	if len(requests) == 0 {
		return nil
	}
	request := `INSERT INTO requests(name, method, target, path) VALUES`
	var requestValues []string
	var values []interface{}
	for _, request := range requests {
		requestValues = append(requestValues, "(?, ?, ?, ?)")
		values = append(values, request.Name, request.Method,
			request.Target.Alias, request.Path)
	}
	request = fmt.Sprintf("%s %s", request, strings.Join(requestValues, ","))

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

func StoreRequest(request models.Request) error {
	return StoreRequests([]models.Request{request})
}

func GetAllRequests() []models.Request {
	request := `SELECT name,method,target,path FROM requests`
	rows, err := globalDB.Query(request)
	if err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return []models.Request{}
	}
	defer rows.Close()

	// Get all targets to fill request.Target
	targetsMap := make(map[string]models.Target)
	for _, target := range GetAllTargets() {
		targetsMap[target.Alias] = target
	}

	var result []models.Request
	for rows.Next() {
		var alias string
		item := models.Request{}
		err := rows.Scan(&item.Name, &item.Method, &alias, &item.Path)
		if err != nil {
			// TODO: log error
			fmt.Printf("error: %+v\n", err)
			return []models.Request{}
		}
		if target, ok := targetsMap[alias]; ok {
			item.Target = target
		} else {
			// TODO: log error
			fmt.Printf("error: alias \"%s\" not found in targets table", alias)
		}
		result = append(result, item)
	}

	return result
}

func GetRequestByName(name string) (models.Request, error) {
	sqlRequest := `SELECT name,method,target,path FROM requests
				WHERE name = ?`
	row := globalDB.QueryRow(sqlRequest, name)
	request := models.Request{}
	if err := row.Scan(&request.Name, &request.Method,
		&request.Target.Alias, &request.Path); err != nil {
		// TODO: log error
		return models.Request{}, err
	}
	target, err := GetTargetByAlias(request.Target.Alias)
	if err != nil {
		// TODO: log error
		return models.Request{}, err
	}
	request.Target = target
	return request, nil
}

func init() {
	if globalDB == nil {
		initDB()
	}
	// create requests table if not exists
	request := `
	CREATE TABLE IF NOT EXISTS requests(
		name TEXT NOT NULL PRIMARY KEY,
		target TEXT,
		method TEXT,
		path TEXT,
		FOREIGN KEY(target) REFERENCES targets(alias)
	);
	`

	_, err := globalDB.Exec(request)
	if err != nil {
		panic(err)
	}
}
