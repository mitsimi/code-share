package domain

type ErrorType string

const (
	ErrNotFound     ErrorType = "NOT_FOUND"
	ErrUnauthorized ErrorType = "UNAUTHORIZED"
	ErrValidation   ErrorType = "VALIDATION"
)

type DomainError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *DomainError {
	return &DomainError{
		Type:    ErrNotFound,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *DomainError {
	return &DomainError{
		Type:    ErrUnauthorized,
		Message: message,
	}
}

func NewValidationError(message string) *DomainError {
	return &DomainError{
		Type:    ErrValidation,
		Message: message,
	}
}
