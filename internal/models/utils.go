package models

import "github.com/mcastorina/poster/internal/store"

func GetRunnableResourceByName(name string) (Runnable, error) {
	resource, err := GetRequestByName(name)
	return &resource, err
}

func GetAllRequests() []Request {
	requests := []Request{}
	for _, sRequest := range store.GetAllRequests() {
		requests = append(requests, convertToRequest(sRequest))
	}
	return requests
}

func GetRequestByName(name string) (Request, error) {
	sRequest, err := store.GetRequestByName(name)
	if err != nil {
		// TODO: log error
		return Request{}, err
	}
	return convertToRequest(sRequest), nil
}

func GetAllEnvironments() []Environment {
	envs := []Environment{}
	for _, sEnvironment := range store.GetAllEnvironments() {
		envs = append(envs, convertToEnvironment(sEnvironment))
	}
	return envs
}
