package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

// Header
type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func (h *Header) String() string {
	return fmt.Sprintf("%s: %s", h.Key, h.Value)
}

// Request
type Request struct {
	Name        string      `yaml:"name"`
	Method      string      `yaml:"method"`
	URL         string      `yaml:"url"`
	Environment Environment `yaml:"environment"`
	Body        string      `yaml:"body"`
	Headers     []Header    `yaml:"headers"`
}

func (r *Request) Run(flags uint32) error {
	return r.RunEnv(r.Environment, flags)
}
func (r *Request) RunEnv(e Environment, flags uint32) error {
	method := e.ReplaceVariables(r.Method)
	url := e.ReplaceVariables(r.URL)
	body := e.ReplaceVariables(r.Body)

	// Create request
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		// TODO: log error
		return err
	}
	// Add headers
	for _, header := range r.Headers {
		key := e.ReplaceVariables(header.Key)
		value := e.ReplaceVariables(header.Value)
		req.Header.Add(key, value)
	}

	// Send request and get response
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
func (r *Request) Save() error {
	if err := r.Validate(); err != nil {
		return err
	}
	return r.ToStore().Save()
}
func (r *Request) Delete() error {
	return r.ToStore().Delete()
}
func (r *Request) Validate() error {
	// Check method is valid
	validMethods := map[string]bool{
		"GET":     true,
		"HEAD":    true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"CONNECT": true,
		"OPTIONS": true,
		"TRACE":   true,
	}
	if _, ok := validMethods[strings.ToUpper(r.Method)]; !ok {
		return ErrorInvalidMethod
	}
	r.Method = strings.ToUpper(r.Method)

	// Check url is valid
	if !strings.Contains(r.URL, "//") {
		r.URL = "//" + r.URL
	}
	urlObj, err := url.Parse(r.URL)
	if err != nil {
		return err
	}
	if urlObj.Scheme == "" {
		urlObj.Scheme = "http"
	}
	r.URL = urlObj.String()
	return nil
}

func flagIsSet(flag, flags uint32) bool {
	return flag&flags != 0
}

// Environment
type Environment struct {
	Name string `yaml:"name"`
}

func (e *Environment) Save() error {
	if err := e.Validate(); err != nil {
		return ErrorInvalidEnvironment
	}
	return e.ToStore().Save()
}
func (e *Environment) Delete() error {
	return e.ToStore().Delete()
}
func (e *Environment) Validate() error {
	return nil
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
	if err := v.Validate(); err != nil {
		return err
	}
	return v.ToStore().Save()
}
func (v *Variable) Delete() error {
	return v.ToStore().Delete()
}
func (v *Variable) Validate() error {
	// Check type is valid
	validTypes := map[string]bool{
		ConstType:   true,
		RequestType: true,
		ScriptType:  true,
	}
	if _, ok := validTypes[strings.ToLower(v.Type)]; !ok {
		return ErrorInvalidType
	}
	v.Type = strings.ToLower(v.Type)
	return nil
}
