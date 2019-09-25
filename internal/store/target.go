package store

import (
	"fmt"
	"strings"

	"github.com/mcastorina/poster/internal/models"
)

func StoreTargets(targets []models.Target) error {
	if len(targets) == 0 {
		return nil
	}
	request := `INSERT INTO targets(alias, url) VALUES`
	var requestValues []string
	var values []interface{}
	for _, target := range targets {
		requestValues = append(requestValues, "(?, ?)")
		values = append(values, target.Alias, target.URL)
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

func StoreTarget(target models.Target) error {
	return StoreTargets([]models.Target{target})
}

func GetAllTargets() []models.Target {
	request := `SELECT alias,url FROM targets`
	rows, err := globalDB.Query(request)
	if err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return []models.Target{}
	}
	defer rows.Close()

	var result []models.Target
	for rows.Next() {
		item := models.Target{}
		err := rows.Scan(&item.Alias, &item.URL)
		if err != nil {
			// TODO: log error
			fmt.Printf("error: %+v\n", err)
			return []models.Target{}
		}
		result = append(result, item)
	}
	return result
}

func GetTargetByAlias(alias string) (models.Target, error) {
	request := `SELECT alias,url FROM targets
				WHERE alias = ?`
	row := globalDB.QueryRow(request, alias)
	target := models.Target{}
	if err := row.Scan(&target.Alias, &target.URL); err != nil {
		// TODO: log error
		return models.Target{}, err
	}
	return target, nil
}

func GetTargetByURL(url string) (models.Target, error) {
	request := `SELECT alias,url FROM targets
				WHERE url = ?`
	row := globalDB.QueryRow(request, url)
	target := models.Target{}
	if err := row.Scan(&target.Alias, &target.URL); err != nil {
		// TODO: log error
		return models.Target{}, err
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
