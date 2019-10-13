package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/mcastorina/poster/internal/cache"
	"github.com/yalp/jsonpath"
)

const (
	ConstType   = "const"
	RequestType = "request"
	ScriptType  = "script"
)

type Resource interface {
	ToStore() interface{}
	Save() error
}

type Runnable interface {
	Run() (*http.Response, error)
	RunEnv(env Environment) (*http.Response, error)
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

// TODO: These private functions are for request variables
//       to avoid infinite loops.
//       Adding a lastGenerated attribute may solve this.
func (r *Request) run() (*http.Response, error) {
	return r.runEnv(r.Environment)
}
func (r *Request) runEnv(e Environment) (*http.Response, error) {
	// Generate variables
	for _, variable := range e.GetVariables() {
		// Skip request variables to avoid infinite loop
		if variable.Type == RequestType {
			continue
		}
		if err := variable.GenerateValue(); err != nil {
			log.Errorf("%+v\n", err)
			return nil, err
		}
		variable.Save()
	}

	method := e.ReplaceVariables(r.Method)
	url := e.ReplaceVariables(r.URL)
	body := e.ReplaceVariables(r.Body)

	// Create request
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Errorf("%+v\n", err)
		return nil, err
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
		log.Errorf("%+v\n", err)
		return nil, err
	}
	return resp, nil
}

func (r *Request) Run() (*http.Response, error) {
	return r.RunEnv(r.Environment)
}
func (r *Request) RunEnv(e Environment) (*http.Response, error) {
	// Generate variables
	for _, variable := range e.GetVariables() {
		if err := variable.GenerateValue(); err != nil {
			log.Errorf("%+v\n", err)
			return nil, errorGenerateVariableFailed
		}
		variable.Save()
	}

	method := e.ReplaceVariables(r.Method)
	url := e.ReplaceVariables(r.URL)
	body := e.ReplaceVariables(r.Body)

	// Create request
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Errorf("%+v\n", err)
		return nil, errorCreateRequestFailed
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
		log.Errorf("%+v\n", err)
		return nil, errorRequestFailed
	}
	return resp, nil
}
func (r *Request) Save() error {
	if err := r.Validate(); err != nil {
		return err
	}
	return cache.SaveRequest(r.ToStore())
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
		return errorInvalidMethod
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

// Environment
type Environment struct {
	Name string `yaml:"name"`
}

func (e *Environment) Save() error {
	if err := e.Validate(); err != nil {
		return errorInvalidEnvironment
	}
	return cache.SaveEnvironment(e.ToStore())
}
func (e *Environment) Delete() error {
	return e.ToStore().Delete()
}
func (e *Environment) Validate() error {
	return nil
}
func (e *Environment) GetVariables() []Variable {
	validVariables := []Variable{}
	for _, variable := range cache.GetVariablesByEnvironment(e.Name) {
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
	Name        string             `yaml:"name"`
	Value       string             `yaml:"value"`
	Environment Environment        `yaml:"environment"`
	Type        string             `yaml:"type"`
	Generator   *VariableGenerator `yaml:"generator,omitempty"`
}

type VariableGenerator struct {
	RequestName string `yaml:"request-name,omitempty"`
	RequestPath string `yaml:"request-path,omitempty"`
	Script      string `yaml:"script,omitempty"`
}

func (v *Variable) Save() error {
	if err := v.Validate(); err != nil {
		return err
	}
	return cache.SaveVariable(v.ToStore())
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
		return errorInvalidType
	}
	v.Type = strings.ToLower(v.Type)
	return nil
}
func (v *Variable) GenerateValue() error {
	switch v.Type {
	case ConstType:
		return nil
	case ScriptType:
		out, err := exec.Command("bash", "-c", v.Generator.Script).Output()
		if err != nil {
			return err
		}
		v.Value = string(bytes.Trim(out, "\n"))
		return nil
	case RequestType:
		req, err := GetRequestByName(v.Generator.RequestName)
		if err != nil {
			return err
		}
		// TODO: add runenv option
		resp, err := req.run()
		if err != nil {
			return err
		}
		// read body
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}

		// assign to variable
		if v.Generator.RequestPath == "" {
			v.Value = string(bytes.Trim(body, "\n"))
		} else {
			var jBody interface{}
			if err := json.Unmarshal(body, &jBody); err != nil {
				return nil
			}
			val, err := jsonpath.Read(jBody, v.Generator.RequestPath)
			if err != nil {
				return err
			}
			if value, ok := val.(string); !ok {
				return fmt.Errorf("the JSONPath points to a non-string value: %+v\n", val)
			} else {
				v.Value = value
			}
		}
		return nil
	}
	return errorInvalidType
}
