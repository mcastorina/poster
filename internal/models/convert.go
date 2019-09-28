package models

import "github.com/mcastorina/poster/internal/store"

func convertToRequest(s store.Request) Request {
	return Request(s)
}
func convertToEnvironment(s store.Environment) Environment {
	return Environment(s)
}
func convertToVariable(s store.Variable) Variable {
	return Variable{
		Name:        s.Name,
		Value:       s.Value,
		Environment: Environment{Name: s.Environment},
		Type:        s.Type,
		Generator:   s.Generator,
	}
}
