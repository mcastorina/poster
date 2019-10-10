package models

import "errors"

var (
	errorInvalidRequest     = errors.New("Request object contains invalid fields")
	errorInvalidEnvironment = errors.New("Environment object contains invalid fields")
	errorInvalidVariable    = errors.New("Variable object contains invalid fields")
	errorInvalidMethod      = errors.New("The provided method is invalid")
	errorInvalidType        = errors.New("The provided type is invalid")

	errorCreateRequestFailed    = errors.New("Could not create a HTTP request")
	errorRequestFailed          = errors.New("Request failed")
	errorGenerateVariableFailed = errors.New("Failed to generate variable")
)
