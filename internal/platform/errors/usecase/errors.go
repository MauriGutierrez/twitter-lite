package usecase

import "fmt"

const (
	TypeNotFound            = "not_found"
	TypeUnknown             = "unknown"
	TypeForbidden           = "forbidden"
	TypeInvalidParam        = "invalid_param"
	TypeInternalServerError = "internal_server_error"
	TypeConflict            = "conflict"
)

type UseCaseError struct {
	Message string
	Type    string
	Causes  []error
}

func (e *UseCaseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func InternalServerError(message string, causes ...error) *UseCaseError {
	return Create(TypeInternalServerError, message, causes...)
}

func NotFound(message string, causes ...error) *UseCaseError {
	return Create(TypeNotFound, message, causes...)
}

func Unknown(message string, causes ...error) *UseCaseError {
	return Create(TypeUnknown, message, causes...)
}

func Forbidden(message string, causes ...error) *UseCaseError {
	return Create(TypeForbidden, message, causes...)
}

func InvalidParam(message string, causes ...error) *UseCaseError {
	return Create(TypeInvalidParam, message, causes...)
}
func Conflict(message string, causes ...error) *UseCaseError {
	return Create(TypeConflict, message, causes...)
}

func Create(Type, message string, errors ...error) *UseCaseError {
	return &UseCaseError{
		Type:    Type,
		Message: message,
		Causes:  errors,
	}
}
