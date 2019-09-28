package store

import "errors"

var (
	ErrorResourceNotFound    = errors.New("resource not found")
	ErrorEnvironmentNotFound = errors.New("environment not found")
	ErrorRequestNotFound     = errors.New("request not found")
	ErrorVariableNotFound    = errors.New("variable not found")
)
