package models

import "errors"

var (
	ErrorInvalidRequest     = errors.New("Request object contains invalid fields")
	ErrorInvalidEnvironment = errors.New("Environment object contains invalid fields")
	ErrorInvalidVariable    = errors.New("Variable object contains invalid fields")
	ErrorInvalidMethod      = errors.New("The provided method is invalid")
	ErrorInvalidType        = errors.New("The provided type is invalid")

	ErrorCreateRequestFailed    = errors.New("Could not create a HTTP request")
	ErrorRequestFailed          = errors.New("Request failed")
	ErrorGenerateVariableFailed = errors.New("Failed to generate variable")
)
