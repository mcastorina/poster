package models

import (
	"github.com/mcastorina/poster/internal/cache"
)

var (
	globalEnvironment = Environment{
		Name: "global",
	}
)

func GetRunnableResourceByName(name string) (Runnable, error) {
	resource, err := GetRequestByName(name)
	return &resource, err
}

func GetAllRequests() []Request {
	requests := []Request{}
	for _, sRequest := range cache.GetAllRequests() {
		requests = append(requests, convertToRequest(sRequest))
	}
	return requests
}
func GetRequestByName(name string) (Request, error) {
	sRequest, err := cache.GetRequestByName(name)
	if err != nil {
		log.Errorf("%+v\n", err)
		return Request{}, err
	}
	return convertToRequest(sRequest), nil
}

func GetAllEnvironments() []Environment {
	envs := []Environment{}
	for _, sEnvironment := range cache.GetAllEnvironments() {
		envs = append(envs, convertToEnvironment(sEnvironment))
	}
	return envs
}
func GetEnvironmentByName(name string) (Environment, error) {
	sEnvironment, err := cache.GetEnvironmentByName(name)
	if err != nil {
		log.Errorf("%+v\n", err)
		return Environment{}, err
	}
	return convertToEnvironment(sEnvironment), nil
}

func GetAllVariables() []Variable {
	vars := []Variable{}
	for _, sVariable := range cache.GetAllVariables() {
		vars = append(vars, convertToVariable(sVariable))
	}
	return vars
}
func GetVariablesByName(name string) []Variable {
	vars := []Variable{}
	for _, sVariable := range cache.GetVariablesByName(name) {
		vars = append(vars, convertToVariable(sVariable))
	}
	return vars
}
func GetVariableByNameAndEnvironment(name, environment string) (Variable, error) {
	sVariable, err := cache.GetVariableByNameAndEnvironment(name, environment)
	if err != nil {
		log.Errorf("%+v\n", err)
		return Variable{}, err
	}
	return convertToVariable(sVariable), nil
}
