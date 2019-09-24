package store

import (
	"fmt"
	"strings"
)

type RequestType struct {
	Name   string
	Method string
	Target TargetType
}

func StoreRequests(requests []RequestType) error {
	if len(requests) == 0 {
		return nil
	}
	request := `INSERT INTO requests(name, method, target) VALUES`
	var requestValues []string
	var values []interface{}
	for _, request := range requests {
		requestValues = append(requestValues, "(?, ?, ?)")
		values = append(values, request.Name, request.Method, request.Target.Alias)
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

func StoreRequest(request RequestType) error {
	return StoreRequests([]RequestType{request})
}

func GetAllRequests() []RequestType {
	request := `SELECT name,method,target FROM requests`
	rows, err := globalDB.Query(request)
	if err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return []RequestType{}
	}
	defer rows.Close()

	// Get all targets to fill request.Target
	targetsMap := make(map[string]TargetType)
	for _, target := range GetAllTargets() {
		targetsMap[target.Alias] = target
	}

	var result []RequestType
	for rows.Next() {
		var alias string
		item := RequestType{}
		err := rows.Scan(&item.Name, &item.Method, &alias)
		if err != nil {
			// TODO: log error
			fmt.Printf("error: %+v\n", err)
			return []RequestType{}
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
