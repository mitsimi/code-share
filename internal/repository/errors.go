package repository

import (
	"errors"
	"fmt"
)

// Common repository errors
var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInternal      = errors.New("internal repository error")
)

// IsNotFound returns true if the error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists returns true if the error is an already exists error
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsInvalidInput returns true if the error is an invalid input error
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// WrapError wraps an error with additional context
func WrapError(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
