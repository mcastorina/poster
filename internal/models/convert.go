package models

import "github.com/mcastorina/poster/internal/store"
import "strings"

func (r *Request) ToStore() *store.Request {
	headerStrings := []string{}
	for _, header := range r.Headers {
		headerStrings = append(headerStrings, header.String())
	}

	return &store.Request{
		Name:        r.Name,
		Method:      r.Method,
		URL:         r.URL,
		Environment: r.Environment.Name,
		Body:        []byte(r.Body),
		Headers:     strings.Join(headerStrings, "\n"),
	}
}
func convertToRequest(s store.Request) Request {
	headers := []Header{}
	if len(s.Headers) > 0 {
		headerStrings := strings.Split(s.Headers, "\n")
		for _, headerString := range headerStrings {
			keyValue := strings.SplitN(headerString, ": ", 2)
			headers = append(headers, Header{Key: keyValue[0], Value: keyValue[1]})
		}
	}
	return Request{
		Name:        s.Name,
		Method:      s.Method,
		URL:         s.URL,
		Environment: Environment{Name: s.Environment},
		Body:        string(s.Body),
		Headers:     headers,
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
