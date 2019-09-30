package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
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
	RunEnv(env Environment, flags uint32) error
}

// Request
type Request struct {
	Name        string      `yaml:"name"`
	Method      string      `yaml:"method"`
	URL         string      `yaml:"url"`
	Environment Environment `yaml:"environment"`
}

func (r *Request) Run(flags uint32) error {
	method := r.Environment.ReplaceVariables(r.Method)
	url := r.Environment.ReplaceVariables(r.URL)

	req, err := http.NewRequest(method, url, nil)
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
func (r *Request) RunEnv(e Environment, flags uint32) error {
	r.Environment = e
	return r.Run(flags)
}
func (r *Request) Save() error {
	return r.ToStore().Save()
}
func (r *Request) Delete() error {
	return r.ToStore().Delete()
}

func flagIsSet(flag, flags uint32) bool {
	return flag&flags != 0
}

// Environment
type Environment struct {
	Name string `yaml:"name"`
}

func (e *Environment) Save() error {
	return e.ToStore().Save()
}
func (e *Environment) Delete() error {
	return e.ToStore().Delete()
}
func (e *Environment) GetVariables() []Variable {
	validVariables := []Variable{}
	for _, variable := range store.GetVariablesByEnvironment(e.Name) {
		validVariables = append(validVariables, convertToVariable(variable))
	}
	return validVariables
}
func (e *Environment) ReplaceVariables(input string) string {
	if strings.Index(input, ":") == -1 {
		return input
	}

	type locType struct {
		startIndex int
		endIndex   int
		value      string
	}
	varLocs := []locType{}

	// Iterate over all variables and create a slice of their locations
	// in the input string.
	for _, variable := range e.GetVariables() {
		re := regexp.MustCompile(`:` + variable.Name + `\b`)
		for _, loc := range re.FindAllStringIndex(input, -1) {
			varLocs = append(varLocs, locType{
				startIndex: loc[0],
				endIndex:   loc[1],
				value:      variable.Value,
			})
		}
	}

	// Reverse sort the list to iterate from largest startIndex to smallest
	sort.Slice(varLocs, func(i, j int) bool {
		return varLocs[i].startIndex > varLocs[j].startIndex
	})

	output := input
	for _, loc := range varLocs {
		// Replace [startIndex:endIndex] with value
		output = output[:loc.startIndex] + loc.value + output[loc.endIndex:]
	}

	return output
}

// Variable
type Variable struct {
	Name        string      `yaml:"name"`
	Value       string      `yaml:"value"`
	Environment Environment `yaml:"environment"`
	Type        string      `yaml:"type"`
	Generator   string      `yaml:"generator"`
}

func (v *Variable) Save() error {
	return v.ToStore().Save()
}
func (v *Variable) Delete() error {
	return v.ToStore().Delete()
}
