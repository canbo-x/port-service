// Package errs contains custom error definitions that are used across the application.
package errs

import "errors"

// Predefined error variables for common errors.
var (
	// ErrPortNotFound is returned when the requested port is not found in the repository.
	ErrPortNotFound = errors.New("port not found")

	// ErrInvalidPortID is returned when the provided port ID is invalid.
	ErrInvalidPortID = errors.New("invalid port id")

	// ErrInvalidInput is returned when the provided input to a function or method is invalid.
	ErrInvalidInput = errors.New("invalid input")
)

// CustomError is a custom error type that can be used for more complex error handling.
type CustomError struct {
	// Message holds the description of the error.
	Message string
}

// Error implements the error interface for CustomError.
func (e *CustomError) Error() string {
	return e.Message
}
