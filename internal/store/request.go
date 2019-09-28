package store

import (
	"fmt"
)

type Request struct {
	Name   string
	Method string
	URL    string
}

func StoreRequests(requests []Request) error {
	if len(requests) == 0 {
		return nil
	}
	tx := globalDB.MustBegin()

	for _, request := range requests {
		tx.NamedExec("INSERT INTO requests (name, method, url) VALUES (:name, :method, :url)",
			&request)
	}

	return tx.Commit()
}

func StoreRequest(request Request) error {
	return StoreRequests([]Request{request})
}

func GetAllRequests() []Request {
	requests := []Request{}
	if err := globalDB.Select(&requests, "SELECT * FROM requests"); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
	}
	return requests
}

func GetRequestByName(name string) (Request, error) {
	request := Request{}
	if err := globalDB.Get(&request, "SELECT * FROM requests WHERE name=$1", name); err != nil {
		// TODO: log error
		fmt.Printf("error: %+v\n", err)
		return Request{}, err
	}
	return request, nil
}

func init() {
	if globalDB == nil {
		initDB()
	}
	// create requests table if not exists
	query := `
	CREATE TABLE IF NOT EXISTS requests(
		name TEXT NOT NULL PRIMARY KEY,
		method TEXT,
		url TEXT
	);
	`

	_, err := globalDB.Exec(query)
	if err != nil {
		panic(err)
	}
}
