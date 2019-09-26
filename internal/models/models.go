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
	Name   string
	Method string
	Target Target
	Path   string
}

func (r *Request) Run() {
	fmt.Printf("%s %s %s %s\n", r.Name, r.Method, r.Target.URL, r.Path)
}
func (r *Request) ToStore() store.Request {
	return store.Request{
		Name:   r.Name,
		Method: r.Method,
		Target: r.Target.Alias,
		Path:   r.Path,
	}
}
func (r *Request) Save() error {
	return store.StoreRequest(r.ToStore())
}

// Target
type Target struct {
	Alias string
	URL   string
}

func (t *Target) ToStore() store.Target {
	return store.Target(*t)
}
func (t *Target) Save() error {
	return store.StoreTarget(t.ToStore())
}

// Environment
type Environment struct {
	Name string
}

func (e *Environment) ToStore() store.Environment {
	return store.Environment(*e)
}
func (e *Environment) Save() error {
	return store.StoreEnvironment(e.ToStore())
}
