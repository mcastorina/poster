package store

import "errors"

var (
	ErrorResourceNotFound    = errors.New("resource not found")
	ErrorEnvironmentNotFound = errors.New("environment not found")
	ErrorEnvironmentInUse    = errors.New("environment is used by other resources")
	ErrorRequestNotFound     = errors.New("request not found")
	ErrorVariableNotFound    = errors.New("variable not found")
	ErrorUnknown             = errors.New("an unknown exception has occurred")
)
