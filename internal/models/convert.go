package models

import (
	"fmt"
	"strings"

	"github.com/mcastorina/poster/internal/store"
)

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
	generator := ""
	switch v.Type {
	case ScriptType:
		generator = v.Generator.Script
	case RequestType:
		generator = fmt.Sprintf("%s:%s", v.Generator.RequestName, v.Generator.RequestPath)
	}
	return &store.Variable{
		Name:        v.Name,
		Value:       v.Value,
		Environment: v.Environment.Name,
		Type:        v.Type,
		Generator:   generator,
	}
}
func convertToVariable(s store.Variable) Variable {
	variable := Variable{
		Name:        s.Name,
		Value:       s.Value,
		Environment: Environment{Name: s.Environment},
		Type:        s.Type,
	}
	generator := &VariableGenerator{}
	switch variable.Type {
	case ScriptType:
		generator.Script = s.Generator
	case RequestType:
		namePath := strings.SplitN(s.Generator, ":", 2)
		generator.RequestName = namePath[0]
		generator.RequestPath = namePath[1]
	case ConstType:
		generator = nil
	}
	variable.Generator = generator
	return variable
}
