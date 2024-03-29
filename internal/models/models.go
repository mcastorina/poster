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
	"time"

	"github.com/mcastorina/poster/internal/cache"
	"github.com/yalp/jsonpath"
)

const (
	ConstType   = "const"
	RequestType = "request"
	ScriptType  = "script"

	variableRegexp = `:([\w-]*)\b`
)

type Resource interface {
	ToStore() interface{}
	Save() error
}

type Runnable interface {
	Run() (*http.Response, error)
	RunEnv(env Environment) (*http.Response, error)
	UpdateHeaders(headers []Header) error
	UpdateBody(body string) error
	UpdateVariables(variables []Variable) error
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

var overrideVariables []Variable

func (r *Request) Run() (*http.Response, error) {
	// TODO: Check for dependency cycles
	return r.RunEnv(r.Environment)
}
func (r *Request) RunEnv(e Environment) (*http.Response, error) {
	// Generate variables
	for _, variable := range e.GetVariablesInRequest(r) {
		if err := variable.GenerateValue(); err != nil {
			log.Errorf("%+v\n", err)
			return nil, errorGenerateVariableFailed
		}
		// TODO: This is a hack to prevent saving override variables
		if variable.Type != ConstType {
			variable.Save()
		}
	}

	methodStr := e.ReplaceVariables(r.Method)
	urlStr := e.ReplaceVariables(r.URL)
	bodyStr := e.ReplaceVariables(r.Body)

	// Check url is valid
	if !strings.Contains(urlStr, "//") {
		urlStr = "//" + urlStr
	}
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		log.Errorf("%+v\n", err)
		return nil, errorInvalidURL
	}
	if urlObj.Scheme == "" {
		urlObj.Scheme = "http"
	}
	urlStr = urlObj.String()

	// Create request
	req, err := http.NewRequest(methodStr, urlStr, strings.NewReader(bodyStr))
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
	// Log sent data
	{
		logMessage := fmt.Sprintf("Sending request:\n> %s %s %s\n", req.Method, req.URL, req.Proto)
		for key, value := range req.Header {
			logMessage += "> " + key + ": " + strings.Join(value, ", ") + "\n"
		}
		logMessage += "\n"
		log.Debugf(logMessage)
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
	return nil
}
func (r *Request) UpdateHeaders(headers []Header) error {
	headerMap := make(map[string]*Header)
	for i, header := range r.Headers {
		headerMap[header.Key] = &r.Headers[i]
	}

	for _, newHeader := range headers {
		if header, ok := headerMap[newHeader.Key]; ok {
			header.Value = newHeader.Value
		} else {
			r.Headers = append(r.Headers, newHeader)
		}
	}
	return nil
}
func (r *Request) UpdateBody(body string) error {
	r.Body = body
	return nil
}
func (r *Request) UpdateVariables(variables []Variable) error {
	overrideVariables = variables
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
func (e *Environment) GetVariableNames() []string {
	varNames := []string{}
	for _, sVariable := range cache.GetVariablesByEnvironment(e.Name) {
		varNames = append(varNames, sVariable.Name)
	}
	sort.Strings(varNames)
	return varNames
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
	// TODO: Optimize this to use only required variables (found in generate step).
	for _, variable := range e.GetVariablesWithGlobal() {
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
func (e *Environment) GetVariablesInRequest(r *Request) []Variable {
	// Map of valid variable names
	validVariables := make(map[string]Variable)
	// Global must be first so it gets overwritten on collision
	for _, variable := range globalEnvironment.GetVariables() {
		validVariables[variable.Name] = variable
	}
	for _, variable := range e.GetVariables() {
		validVariables[variable.Name] = variable
	}
	for _, variable := range overrideVariables {
		validVariables[variable.Name] = variable
	}
	// Build search string as a combination of all parts
	// of the request that can be replaced
	searchString := r.Method + "\n" + r.URL + "\n" + r.Body
	for _, header := range r.Headers {
		searchString = searchString + "\n" + header.Key + "\n" + header.Value
	}

	// Search for variables in the string and add to slice
	// if it is a valid variable name
	variables := []Variable{}
	re := regexp.MustCompile(variableRegexp)
	for _, varGroup := range re.FindAllStringSubmatch(searchString, -1) {
		varName := varGroup[1]
		if variable, ok := validVariables[varName]; ok {
			variables = append(variables, variable)
		}
	}
	return variables
}
func (e *Environment) GetVariablesWithGlobal() []Variable {
	// Map of valid variable names
	validVariables := make(map[string]Variable)
	// Global must be first so it gets overwritten on collision
	for _, variable := range globalEnvironment.GetVariables() {
		validVariables[variable.Name] = variable
	}
	for _, variable := range e.GetVariables() {
		validVariables[variable.Name] = variable
	}
	for _, variable := range overrideVariables {
		validVariables[variable.Name] = variable
	}
	// Build return array
	variables := []Variable{}
	for _, variable := range validVariables {
		variables = append(variables, variable)
	}
	return variables
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
	RequestName        string `yaml:"request-name,omitempty"`
	RequestPath        string `yaml:"request-path,omitempty"`
	RequestEnvironment string `yaml:"request-environment,omitempty"`
	Script             string `yaml:"script,omitempty"`
	Timeout            int64  `yaml:"timeout"`
	LastGenerated      time.Time
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
	// Check that the regexp matches
	re := regexp.MustCompile("^" + variableRegexp + "$")
	if !re.MatchString(":" + v.Name) {
		return errorInvalidCharacters
	}
	// TODO: Verify generator
	return nil
}
func (v *Variable) GenerateValue() error {
	if v.Generator == nil {
		return nil
	}
	timeout := time.Duration(v.Generator.Timeout) * time.Minute
	if time.Since(v.Generator.LastGenerated) < timeout {
		return nil
	}
	log.Infof("Variable %s is stale, generating new value..", v.Name)

	switch v.Type {
	case ConstType:
		return nil
	case ScriptType:
		out, err := exec.Command("bash", "-c", v.Generator.Script).Output()
		if err != nil {
			return err
		}
		v.Value = string(bytes.Trim(out, "\n"))
		v.Generator.LastGenerated = time.Now()
		return nil
	case RequestType:
		req, err := GetRequestByName(v.Generator.RequestName)
		if err != nil {
			return err
		}
		env, err := GetEnvironmentByName(v.Generator.RequestEnvironment)
		if err != nil {
			return err
		}
		resp, err := req.RunEnv(env)
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
				return err
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
		log.Debugf("Variable %s updated to: %s\n", v.Name, v.Value)
		v.Generator.LastGenerated = time.Now()
		return nil
	}
	return errorInvalidType
}
