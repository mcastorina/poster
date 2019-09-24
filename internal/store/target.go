package store

import (
	"fmt"
	"strings"
)

type TargetType struct {
	Alias string
	URL   string
}

func StoreTargets(targets []TargetType) error {
	if len(targets) == 0 {
		return nil
	}
	request := `
	INSERT INTO targets(
		alias,
		url
	) VALUES
	`
	var requestValues []string
	var values []interface{}
	for _, target := range targets {
		requestValues = append(requestValues, "(?, ?)")
		values = append(values, target.Alias, target.URL)
	}
	request = fmt.Sprintf("%s %s;", request, strings.Join(requestValues, ","))

	stmt, err := globalDB.Prepare(request)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return err
	}
	return nil
}

func StoreTarget(target TargetType) error {
	return StoreTargets([]TargetType{target})
}

func GetAllTargets() []TargetType {
	request := `
	SELECT alias,url FROM targets;
	`
	rows, err := globalDB.Query(request)
	if err != nil {
		fmt.Printf("error: %+v\n", err)
		return []TargetType{}
	}
	defer rows.Close()

	var result []TargetType
	for rows.Next() {
		item := TargetType{}
		err := rows.Scan(&item.Alias, &item.URL)
		if err != nil {
			fmt.Printf("error: %+v\n", err)
			return []TargetType{}
		}
		result = append(result, item)
	}
	return result
}
