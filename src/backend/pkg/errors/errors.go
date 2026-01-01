package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error types for domain errors
type ErrorType string

const (
	ErrorTypeNotFound       ErrorType = "NOT_FOUND"
	ErrorTypeValidation     ErrorType = "VALIDATION"
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden      ErrorType = "FORBIDDEN"
	ErrorTypeConflict       ErrorType = "CONFLICT"
	ErrorTypeInternal       ErrorType = "INTERNAL"
	ErrorTypeBadRequest     ErrorType = "BAD_REQUEST"
	ErrorTypeServiceUnavail ErrorType = "SERVICE_UNAVAILABLE"
)

// AppError represents an application error
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// GRPCStatus converts the error to a gRPC status
func (e *AppError) GRPCStatus() *status.Status {
	code := codes.Internal
	switch e.Type {
	case ErrorTypeNotFound:
		code = codes.NotFound
	case ErrorTypeValidation, ErrorTypeBadRequest:
		code = codes.InvalidArgument
	case ErrorTypeUnauthorized:
		code = codes.Unauthenticated
	case ErrorTypeForbidden:
		code = codes.PermissionDenied
	case ErrorTypeConflict:
		code = codes.AlreadyExists
	case ErrorTypeServiceUnavail:
		code = codes.Unavailable
	}
	return status.New(code, e.Message)
}

// Constructor functions

// NotFound creates a not found error
func NotFound(resource string, id string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s with ID %s not found", resource, id),
	}
}

// Validation creates a validation error
func Validation(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
	}
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
	}
}

// Forbidden creates a forbidden error
func Forbidden(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
	}
}

// Conflict creates a conflict error
func Conflict(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
	}
}

// Internal creates an internal error
func Internal(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Err:     err,
	}
}

// BadRequest creates a bad request error
func BadRequest(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
	}
}

// ServiceUnavailable creates a service unavailable error
func ServiceUnavailable(service string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeServiceUnavail,
		Message: fmt.Sprintf("%s service is unavailable", service),
		Err:     err,
	}
}

// ToGRPCError converts any error to a gRPC status error
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	if appErr, ok := err.(*AppError); ok {
		return appErr.GRPCStatus().Err()
	}

	return status.Error(codes.Internal, "internal server error")
}
