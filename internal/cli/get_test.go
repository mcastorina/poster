package cli

import (
	"testing"

	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestsFromArgumentsAll(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	requests := getRequestsFromArguments("", "", emptyArr, emptyArr, emptyArr, emptyArr)

	assert.Equal(t, 4, len(requests))
}
func TestGetRequestsFromArgumentsWithEnvFlag(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	requests := getRequestsFromArguments("test", "", emptyArr, emptyArr, emptyArr, emptyArr)
	assert.Equal(t, 4, len(requests))

	requests = getRequestsFromArguments("remote", "", emptyArr, emptyArr, emptyArr, emptyArr)
	assert.Equal(t, 0, len(requests))
}
func TestGetRequestsFromArgumentsWithMethodFlag(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	requests := getRequestsFromArguments("", "GET", emptyArr, emptyArr, emptyArr, emptyArr)
	assert.Equal(t, 2, len(requests))
	assert.Equal(t, "test1", requests[0].Name)
	assert.Equal(t, "test2", requests[1].Name)
}
func TestGetRequestsFromArgumentsWithVariables(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	requests := getRequestsFromArguments("", "", []string{"host"}, emptyArr, emptyArr, emptyArr)
	assert.Equal(t, 3, len(requests))

	requests = getRequestsFromArguments("", "", []string{"token"}, emptyArr, emptyArr, emptyArr)
	assert.Equal(t, 2, len(requests))

	requests = getRequestsFromArguments("", "", []string{"token", "host"}, emptyArr, emptyArr, emptyArr)
	assert.Equal(t, 1, len(requests))
	assert.Equal(t, "test2", requests[0].Name)
}
func TestGetRequestsFromArgumentsWithHeaders(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	requests := getRequestsFromArguments("", "", emptyArr, []string{"application/json"}, emptyArr, emptyArr)
	assert.Equal(t, 1, len(requests))
	assert.Equal(t, "test3", requests[0].Name)

	requests = getRequestsFromArguments("", "", emptyArr, []string{"Authorization"}, emptyArr, emptyArr)
	assert.Equal(t, 2, len(requests))
}
func TestGetRequestsFromArgumentsWithBody(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	requests := getRequestsFromArguments("", "", emptyArr, emptyArr, []string{"hello"}, emptyArr)
	assert.Equal(t, 2, len(requests))
	if requests[0].Name == "test3" {
		assert.Equal(t, "test3", requests[0].Name)
		assert.Equal(t, "test4", requests[1].Name)
	} else {
		assert.Equal(t, "test4", requests[0].Name)
		assert.Equal(t, "test3", requests[1].Name)
	}

	requests = getRequestsFromArguments("", "", emptyArr, emptyArr, []string{"hello", "foo"}, emptyArr)
	assert.Equal(t, 0, len(requests))

	requests = getRequestsFromArguments("", "", emptyArr, emptyArr, []string{"not found"}, emptyArr)
	assert.Equal(t, 0, len(requests))
}
func TestGetRequestsFromArgumentsWithName(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	args := []string{"test1", "test3"}
	requests := getRequestsFromArguments("", "", emptyArr, emptyArr, emptyArr, args)

	assert.Equal(t, 2, len(requests))
	assert.Equal(t, "test1", requests[0].Name)
	assert.Equal(t, "test3", requests[1].Name)
}
func TestGetRequestsFromArgumentsWithEnvFlagAndMethodFlagAndVariables(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	requests := getRequestsFromArguments("test", "POST", []string{"host"}, emptyArr, emptyArr, emptyArr)

	assert.Equal(t, 1, len(requests))
	assert.Equal(t, "test3", requests[0].Name)
}
func TestGetRequestsFromArgumentsWithMethodFlagAndVariablesAndHeaders(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	requests := getRequestsFromArguments("", "POST", []string{"host"}, []string{"application/json"}, emptyArr, emptyArr)

	assert.Equal(t, 1, len(requests))
	assert.Equal(t, "test3", requests[0].Name)

	requests = getRequestsFromArguments("", "POST", []string{"host"}, []string{"Authorization"}, emptyArr, emptyArr)
	assert.Equal(t, 0, len(requests))
}
func TestGetRequestsFromArgumentsWithMethodFlagAndVariablesAndHeadersAndName(t *testing.T) {
	defer monkey.UnpatchAll()
	patchGetRequests()

	emptyArr := []string{}
	requests := getRequestsFromArguments("", "POST", []string{"host"}, []string{"application/json"}, emptyArr, []string{"test2", "test3"})

	assert.Equal(t, 1, len(requests))
	assert.Equal(t, "test3", requests[0].Name)

	requests = getRequestsFromArguments("", "GET", []string{"host", "token"}, []string{"Authorization"}, emptyArr, []string{"test2", "test3"})
	assert.Equal(t, 1, len(requests))
	assert.Equal(t, "test2", requests[0].Name)

	requests = getRequestsFromArguments("", "GET", []string{"host", "token"}, []string{"Authorization"}, emptyArr, []string{"test1", "test3"})
	assert.Equal(t, 0, len(requests))
}
