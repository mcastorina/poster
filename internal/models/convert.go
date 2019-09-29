package models

import "github.com/mcastorina/poster/internal/store"

func (r *Request) ToStore() *store.Request {
	return &store.Request{
		Name:        r.Name,
		Method:      r.Method,
		URL:         r.URL,
		Environment: r.Environment.Name,
	}
}
func convertToRequest(s store.Request) Request {
	return Request{
		Name:        s.Name,
		Method:      s.Method,
		URL:         s.URL,
		Environment: Environment{Name: s.Environment},
	}
}

func (e *Environment) ToStore() *store.Environment {
	sEnv := store.Environment(*e)
	return &sEnv
}
func convertToEnvironment(s store.Environment) Environment {
	return Environment(s)
}

func (v *Variable) ToStore() *store.Variable {
	return &store.Variable{
		Name:        v.Name,
		Value:       v.Value,
		Environment: v.Environment.Name,
		Type:        v.Type,
		Generator:   v.Generator,
	}
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
