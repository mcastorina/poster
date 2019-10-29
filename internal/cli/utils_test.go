package cli

import (
	"fmt"
	"reflect"

	"github.com/bouk/monkey"
	"github.com/mcastorina/poster/internal/models"
)

var testEnvironment = models.Environment{
	Name: "test",
}
var testRequests = []models.Request{
	{
		Name:        "test1",
		Method:      "GET",
		URL:         ":host",
		Environment: testEnvironment,
		Body:        "",
		Headers:     []models.Header{},
	},
	{
		Name:        "test2",
		Method:      "GET",
		URL:         ":host",
		Environment: testEnvironment,
		Body:        "",
		Headers: []models.Header{
			{
				Key:   "Authorization",
				Value: ":token",
			},
		},
	},
	{
		Name:        "test3",
		Method:      "POST",
		URL:         ":host",
		Environment: testEnvironment,
		Body:        `{"hello": "world"}`,
		Headers: []models.Header{
			{
				Key:   "Content-Type",
				Value: "application/json",
			},
		},
	},
	{
		Name:        "test4",
		Method:      "POST",
		URL:         "localhost",
		Environment: testEnvironment,
		Body:        `{"hello": "world"}`,
		Headers: []models.Header{
			{
				Key:   "Authorization",
				Value: ":token",
			},
		},
	},
}

func patchGetRequests() {
	// models.GetRequestsByEnvironmentAndMethod(envFlag, methodFlag)
	patchGetRequestsByEnvironmentAndMethod()
	// models.GetRequestByName(arg)
	patchGetRequestByName()
	// models.GetRequestsByEnvironment(envFlag)
	patchGetRequestsByEnvironment()
	// models.GetRequestsByMethod(methodFlag)
	patchGetRequestsByMethod()
	// models.GetAllRequests()
	patchGetAllRequests()
	// request.Environment.GetVariablesInRequest(&request)
	patchGetVariablesInRequest()
}

func patchGetVariablesInRequest() {
	patch := func(e *models.Environment, r *models.Request) []models.Variable {
		if e.Name != testEnvironment.Name {
			return []models.Variable{}
		}
		hostVar := models.Variable{Name: "host", Type: models.ConstType}
		tokenVar := models.Variable{Name: "token", Type: models.RequestType}
		switch r.Name {
		case "test1":
			return []models.Variable{hostVar}
		case "test2":
			return []models.Variable{hostVar, tokenVar}
		case "test3":
			return []models.Variable{hostVar}
		case "test4":
			return []models.Variable{tokenVar}
		}
		return []models.Variable{}
	}
	var envPtr *models.Environment
	monkey.PatchInstanceMethod(reflect.TypeOf(envPtr), "GetVariablesInRequest", patch)
}
func patchGetRequestsByEnvironmentAndMethod() {
	patch := func(env, method string) []models.Request {
		switch method {
		case "GET":
			return testRequests[0:2]
		case "POST":
			return testRequests[2:4]
		}
		return []models.Request{}
	}
	monkey.Patch(models.GetRequestsByEnvironmentAndMethod, patch)
}
func patchGetRequestsByEnvironment() {
	patch := func(env string) []models.Request {
		if env == testEnvironment.Name {
			return testRequests
		}
		return []models.Request{}
	}
	monkey.Patch(models.GetRequestsByEnvironment, patch)
}
func patchGetRequestsByMethod() {
	patch := func(method string) []models.Request {
		switch method {
		case "GET":
			return testRequests[0:2]
		case "POST":
			return testRequests[2:4]
		}
		return []models.Request{}
	}
	monkey.Patch(models.GetRequestsByMethod, patch)
}
func patchGetAllRequests() {
	patch := func() []models.Request {
		return testRequests
	}
	monkey.Patch(models.GetAllRequests, patch)
}
func patchGetRequestByName() {
	patch := func(name string) (models.Request, error) {
		switch name {
		case "test1":
			return testRequests[0], nil
		case "test2":
			return testRequests[1], nil
		case "test3":
			return testRequests[2], nil
		case "test4":
			return testRequests[3], nil
		}
		return models.Request{}, fmt.Errorf("not found")
	}
	monkey.Patch(models.GetRequestByName, patch)
}
