package models

import "github.com/mcastorina/poster/internal/store"

func convertToRequest(s store.Request) Request {
	sTarget, err := store.GetTargetByAlias(s.Target)
	if err != nil {
		// TODO: log error
		return Request{}
	}
	r := Request{
		Name:   s.Name,
		Method: s.Method,
		Target: convertToTarget(sTarget),
		Path:   s.Path,
	}

	return r
}
func convertToTarget(s store.Target) Target {
	return Target(s)
}
func convertToEnvironment(s store.Environment) Environment {
	return Environment(s)
}
