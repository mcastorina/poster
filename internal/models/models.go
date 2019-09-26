package models

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mcastorina/poster/internal/store"
)

type Resource interface {
	ToStore() interface{}
	Save() error
}

type Runnable interface {
	Run() error
}

// Request
type Request struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

func (r *Request) Run() error {
	req, err := http.NewRequest(r.Method, r.URL, nil)
	if err != nil {
		// TODO: log error
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// TODO: log error
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	fmt.Printf("%s", body)
	return nil
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
