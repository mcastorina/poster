package models

import (
	"fmt"
	"github.com/mcastorina/poster/internal/store"
)

type Resource interface {
	ToStore() interface{}
	Save() error
}

type Runnable interface {
	Run()
}

// Request
type Request struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

func (r *Request) Run() {
	fmt.Printf("%s %s %s\n", r.Name, r.Method, r.URL)
}
func (r *Request) ToStore() store.Request {
	return store.Request(*r)
}
func (r *Request) Save() error {
	return store.StoreRequest(r.ToStore())
}

// Environment
type Environment struct {
	Name string `json:"name"`
}

func (e *Environment) ToStore() store.Environment {
	return store.Environment(*e)
}
func (e *Environment) Save() error {
	return store.StoreEnvironment(e.ToStore())
}
