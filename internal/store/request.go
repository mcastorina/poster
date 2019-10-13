package store

import (
	"github.com/mattn/go-sqlite3"
)

type Request struct {
	Name        string
	Method      string
	URL         string
	Environment string
	Body        []byte
	Headers     string // newline separated values
}

func (r *Request) Save() error {
	return StoreRequests([]Request{*r})
}
func (r *Request) Delete() error {
	_, err := globalDB.Exec("DELETE FROM requests WHERE name=$1", r.Name)
	if err != nil {
		log.Errorf("%+v\n", err)
		return ErrorRequestNotFound
	}
	return nil
}

func StoreRequests(requests []Request) error {
	if len(requests) == 0 {
		return nil
	}
	tx := globalDB.MustBegin()

	for _, request := range requests {
		if _, err := tx.NamedExec(
			`INSERT OR REPLACE INTO requests
			(name, method, url, environment, body, headers)
			VALUES (:name, :method, :url, :environment, :body, :headers)`,
			&request); err != nil {

			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.Code == sqlite3.ErrConstraint {
					if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
						return ErrorRequestExists
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

func GetAllRequests() []Request {
	requests := []Request{}
	if err := globalDB.Select(&requests, "SELECT * FROM requests"); err != nil {
		log.Errorf("%+v\n", err)
	}
	return requests
}
func GetRequestByName(name string) (Request, error) {
	request := Request{}
	if err := globalDB.Get(&request, "SELECT * FROM requests WHERE name=$1", name); err != nil {
		log.Errorf("%+v\n", err)
		return Request{}, ErrorRequestNotFound
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
		url TEXT,
		environment TEXT,
		body BLOB,
		headers TEXT,
		FOREIGN KEY(environment) REFERENCES environments(name)
	);
	`

	_, err := globalDB.Exec(query)
	if err != nil {
		panic(err)
	}
}
