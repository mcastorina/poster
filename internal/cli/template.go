package cli

import (
	"github.com/mcastorina/poster/internal/models"
)

type Request struct {
	Name        string            `yaml:"name"`
	Method      string            `yaml:"method"`
	URL         string            `yaml:"url"`
	Environment string            `yaml:"default-environment"`
	Body        string            `yaml:"body,omitempty"`
	Headers     map[string]string `yaml:"headers"`
}

func (r *Request) Save() error {
	env, err := models.GetEnvironmentByName(r.Environment)
	if err != nil {
		return err
	}
	headers := []models.Header{}
	for key, value := range r.Headers {
		headers = append(headers, models.Header{
			Key:   key,
			Value: value,
		})
	}
	request := models.Request{
		Name:        r.Name,
		Method:      r.Method,
		URL:         r.URL,
		Environment: env,
		Body:        r.Body,
		Headers:     headers,
	}
	return request.Save()
}

type Environment struct {
	Name      string   `yaml:"name"`
	Variables []string `yaml:"variables"`
}

func (e *Environment) Save() error {
	env := models.Environment{
		Name: e.Name,
	}
	// TODO: Do something with e.Variables
	return env.Save()
}

type Variable struct {
	Name         string             `yaml:"name"`
	Type         string             `yaml:"type"`
	Value        string             `yaml:"value,omitempty"`
	Environments []string           `yaml:"environments"`
	Generator    *VariableGenerator `yaml:"generator,omitempty"`
}
type VariableGenerator struct {
	RequestName        string `yaml:"name,omitempty"`
	RequestPath        string `yaml:"jsonpath,omitempty"`
	RequestEnvironment string `yaml:"environment,omitempty"`
	Script             string `yaml:"script,omitempty"`
}

func (v *Variable) Save() error {
	var generator *models.VariableGenerator
	parent := false
	if v.Generator != nil {
		generatorStruct := models.VariableGenerator(*v.Generator)
		generator = &generatorStruct
		parent = generator.RequestEnvironment == "parent"
	}
	variable := models.Variable{
		Name:      v.Name,
		Value:     v.Value,
		Type:      v.Type,
		Generator: generator,
	}
	for _, environment := range v.Environments {
		env, err := models.GetEnvironmentByName(environment)
		if err != nil {
			return err
		}
		variable.Environment = env
		if parent {
			generator.RequestEnvironment = env.Name
		}
		if err := variable.Save(); err != nil {
			return err
		}
	}
	return nil
}

func requestTemplate() Request {
	return Request{
		Name:        "my-super-awesome-request",
		Method:      "GET",
		URL:         "http://localhost",
		Environment: "local",
		Body:        `{"msg": "why does a GET request have a body?"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func environmentTemplate() Environment {
	return Environment{
		Name:      "my-super-awesome-environment",
		Variables: []string{},
	}
}

func constVariableTemplate() Variable {
	return Variable{
		Name:         "my-super-awesome-variable",
		Value:        "value",
		Environments: []string{"global"},
		Type:         models.ConstType,
	}
}

func scriptVariableTemplate() Variable {
	return Variable{
		Name:         "my-super-awesome-variable",
		Type:         models.ScriptType,
		Environments: []string{"global"},
		Generator: &VariableGenerator{
			Script: `date +'%D %T'`,
		},
	}
}

func requestVariableTemplate() Variable {
	return Variable{
		Name:         "my-super-awesome-variable",
		Type:         models.RequestType,
		Environments: []string{"global"},
		Generator: &VariableGenerator{
			RequestName:        "request-name",
			RequestEnvironment: "parent",
			RequestPath:        "$",
		},
	}
}
