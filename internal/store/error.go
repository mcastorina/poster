package store

import "errors"

var (
	ErrorResourceNotFound    = errors.New("resource not found")
	ErrorEnvironmentNotFound = errors.New("environment not found")
	ErrorEnvironmentInUse    = errors.New("environment is used by other resources")
	ErrorEnvironmentExists   = errors.New("environment already exists")
	ErrorRequestNotFound     = errors.New("request not found")
	ErrorRequestExists       = errors.New("request already exists")
	ErrorVariableNotFound    = errors.New("variable not found")
	ErrorVariableExists      = errors.New("variable already exists")
	ErrorUnknown             = errors.New("an unknown exception has occurred")
)
