package store

import (
	"fmt"
)

type Request struct {
	Name   string
	Target string
	Method string
	Path   string
}

func StoreRequests(requests []Request) error {
	if len(requests) == 0 {
		return nil
	}
	tx := globalDB.MustBegin()

	for _, request := range requests {
		tx.NamedExec("INSERT INTO requests (name, target, method, path) VALUES (:name, :target, :method, :path)",
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
