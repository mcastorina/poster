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
	sVariable := &store.Variable{
		Name:        v.Name,
		Value:       v.Value,
		Environment: v.Environment.Name,
		Type:        v.Type,
	}
	if v.Generator != nil {
		switch v.Type {
		case ScriptType:
			sVariable.Generator = v.Generator.Script
		case RequestType:
			sVariable.Generator = fmt.Sprintf("%s:%s:%s",
				v.Generator.RequestName, v.Generator.RequestEnvironment,
				v.Generator.RequestPath)
		}
		sVariable.Timeout = v.Generator.Timeout
		sVariable.Last = v.Generator.LastGenerated
	}
	return sVariable
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
		namePath := strings.SplitN(s.Generator, ":", 3)
		generator.RequestName = namePath[0]
		generator.RequestEnvironment = namePath[1]
		generator.RequestPath = namePath[2]
	case ConstType:
		generator = nil
	}
	if generator != nil {
		generator.Timeout = s.Timeout
		generator.LastGenerated = s.Last
	}
	variable.Generator = generator
	return variable
}
