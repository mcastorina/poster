package models

import "errors"

var (
	ErrorInvalidRequest     = errors.New("request object contains invalid fields")
	ErrorInvalidEnvironment = errors.New("environment object contains invalid fields")
	ErrorInvalidVariable    = errors.New("variable object contains invalid fields")
	ErrorInvalidMethod      = errors.New("the provided method is invalid")
	ErrorInvalidType        = errors.New("the provided type is invalid")
)
