package models

import "github.com/mcastorina/poster/internal/store"

func convertToRequest(s store.Request) Request {
	return Request(s)
}
func convertToEnvironment(s store.Environment) Environment {
	return Environment(s)
}
