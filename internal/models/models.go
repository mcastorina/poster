package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mcastorina/poster/internal/store"
)

const (
	FlagPrintResponseCode = uint32(0x1 << 0)
	FlagPrintHeaders      = uint32(0x1 << 1)
	FlagPrintBody         = uint32(0x1 << 2)

	ConstType   = "const"
	RequestType = "request"
	ScriptType  = "script"
)

type Resource interface {
	ToStore() interface{}
	Save() error
}

type Runnable interface {
	Run(flags uint32) error
}

// Request
type Request struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

func (r *Request) Run(flags uint32) error {
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

	if flagIsSet(FlagPrintResponseCode, flags) {
		fmt.Printf("%s %s\n", resp.Proto, resp.Status)
	}
	if flagIsSet(FlagPrintHeaders, flags) {
		for header, values := range resp.Header {
			fmt.Printf("%s: %s\n", header, strings.Join(values, ","))
		}
		fmt.Println()
	}
	if flagIsSet(FlagPrintBody, flags) {
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
		fmt.Printf("%s", body)
	}
	return nil
}
func (r *Request) ToStore() store.Request {
	return store.Request(*r)
}
func (r *Request) Save() error {
	return store.StoreRequest(r.ToStore())
}

func flagIsSet(flag, flags uint32) bool {
	return flag&flags != 0
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

// Variable
type Variable struct {
	Name        string
	Value       string
	Environment Environment
	Type        string
	Generator   string
}

func (v *Variable) ToStore() store.Variable {
	return store.Variable{
		Name:        v.Name,
		Value:       v.Value,
		Environment: v.Environment.Name,
		Type:        v.Type,
		Generator:   v.Generator,
	}
}
func (v *Variable) Save() error {
	return store.StoreVariable(v.ToStore())
}
