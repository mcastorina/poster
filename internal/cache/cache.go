package cache

import (
	"github.com/mcastorina/poster/internal/store"
)

var cache map[string]interface{}

func GetAllRequests() []store.Request {
	key := "GetAllRequests"
	if requests, ok := cacheGet(key); ok {
		return requests.([]store.Request)
	}
	requests := store.GetAllRequests()
	cacheSet(key, requests)
	return requests
}
func GetRequestByName(name string) (store.Request, error) {
	key := "GetRequestByName:" + name
	if result, ok := cacheGet(key); ok {
		return result.(store.Request), nil
	}
	request, err := store.GetRequestByName(name)
	if err != nil {
		return store.Request{}, err
	}
	cacheSet(key, request)
	return request, nil
}
func SaveRequest(r *store.Request) error {
	delete(cache, "GetAllRequests")
	delete(cache, "GetRequestByName:"+r.Name)
	return r.Save()
}

func GetAllEnvironments() []store.Environment {
	key := "GetAllEnvironments"
	if environments, ok := cacheGet(key); ok {
		return environments.([]store.Environment)
	}
	environments := store.GetAllEnvironments()
	cacheSet(key, environments)
	return environments
}
func GetEnvironmentByName(name string) (store.Environment, error) {
	key := "GetEnvironmentByName:" + name
	if result, ok := cacheGet(key); ok {
		return result.(store.Environment), nil
	}
	environment, err := store.GetEnvironmentByName(name)
	if err != nil {
		return store.Environment{}, err
	}
	cacheSet(key, environment)
	return environment, nil
}
func SaveEnvironment(e *store.Environment) error {
	delete(cache, "GetAllEnvironments")
	delete(cache, "GetEnvironmentByName:"+e.Name)
	return e.Save()
}

func GetAllVariables() []store.Variable {
	key := "GetAllVariables"
	if variables, ok := cacheGet(key); ok {
		return variables.([]store.Variable)
	}
	variables := store.GetAllVariables()
	cacheSet(key, variables)
	return variables
}
func GetVariablesByName(name string) []store.Variable {
	key := "GetVariablesByName:" + name
	if variables, ok := cacheGet(key); ok {
		return variables.([]store.Variable)
	}
	variables := store.GetVariablesByName(name)
	cacheSet(key, variables)
	return variables
}
func GetVariablesByEnvironment(environment string) []store.Variable {
	key := "GetVariablesByEnvironment:" + environment
	if variables, ok := cacheGet(key); ok {
		return variables.([]store.Variable)
	}
	variables := store.GetVariablesByEnvironment(environment)
	cacheSet(key, variables)
	return variables
}
func GetVariablesByType(typ string) []store.Variable {
	key := "GetVariablesByType:" + typ
	if variables, ok := cacheGet(key); ok {
		return variables.([]store.Variable)
	}
	variables := store.GetVariablesByType(typ)
	cacheSet(key, variables)
	return variables
}
func GetVariableByNameAndEnvironment(name, environment string) (store.Variable, error) {
	// This is susceptible to collisions but is fine for now
	key := "GetVariableByNameAndEnvironment:" + name + "," + environment
	if variable, ok := cacheGet(key); ok {
		return variable.(store.Variable), nil
	}
	variable, err := store.GetVariableByNameAndEnvironment(name, environment)
	if err != nil {
		return store.Variable{}, err
	}
	cacheSet(key, variable)
	return variable, nil
}
func GetVariablesByNameAndType(name, typ string) []store.Variable {
	key := "GetVariablesByNameAndType:" + name + "," + typ
	if variables, ok := cacheGet(key); ok {
		return variables.([]store.Variable)
	}
	variables := store.GetVariablesByNameAndType(name, typ)
	cacheSet(key, variables)
	return variables
}
func GetVariablesByEnvironmentAndType(environment, typ string) []store.Variable {
	// This is susceptible to collisions but is fine for now
	key := "GetVariablesByEnvironmentAndType:" + environment + "," + typ
	if variables, ok := cacheGet(key); ok {
		return variables.([]store.Variable)
	}
	variables := store.GetVariablesByEnvironmentAndType(environment, typ)
	cacheSet(key, variables)
	return variables
}
func SaveVariable(v *store.Variable) error {
	delete(cache, "GetAllVariables")
	delete(cache, "GetVariablesByEnvironment:"+v.Environment)
	delete(cache, "GetVariablesByName:"+v.Name)
	delete(cache, "GetVariableByNameAndEnvironment:"+v.Name+","+v.Environment)
	return v.Save()
}

func cacheGet(key string) (interface{}, bool) {
	if value, ok := cache[key]; ok {
		log.Debugf("Cache hit on [%s]", key)
		return value, true
	}
	log.Debugf("Cache miss on [%s]", key)
	return nil, false
}
func cacheSet(key string, value interface{}) {
	cache[key] = value
}

func init() {
	cache = make(map[string]interface{})
}
