package tasmota

import (
	"errors"
	"fmt"
)

// ErrorType represents different categories of errors that can occur.
type ErrorType int

const (
	// ErrorTypeNetwork indicates a network-related error.
	ErrorTypeNetwork ErrorType = iota
	// ErrorTypeAuth indicates an authentication error.
	ErrorTypeAuth
	// ErrorTypeCommand indicates a command execution error.
	ErrorTypeCommand
	// ErrorTypeParse indicates a response parsing error.
	ErrorTypeParse
	// ErrorTypeTimeout indicates a timeout error.
	ErrorTypeTimeout
	// ErrorTypeDevice indicates a device-specific error.
	ErrorTypeDevice
)

// String returns a string representation of the ErrorType.
func (e ErrorType) String() string {
	switch e {
	case ErrorTypeNetwork:
		return "network"
	case ErrorTypeAuth:
		return "auth"
	case ErrorTypeCommand:
		return "command"
	case ErrorTypeParse:
		return "parse"
	case ErrorTypeTimeout:
		return "timeout"
	case ErrorTypeDevice:
		return "device"
	default:
		return "unknown"
	}
}

// Error represents a Tasmota client error with additional context.
type Error struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s error: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new Error with the given type and message.
func NewError(errType ErrorType, message string, err error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// IsNetworkError checks if the error is a network error.
func IsNetworkError(err error) bool {
	var tasErr *Error
	if errors.As(err, &tasErr) {
		return tasErr.Type == ErrorTypeNetwork
	}
	return false
}

// IsAuthError checks if the error is an authentication error.
func IsAuthError(err error) bool {
	var tasErr *Error
	if errors.As(err, &tasErr) {
		return tasErr.Type == ErrorTypeAuth
	}
	return false
}

// IsCommandError checks if the error is a command execution error.
func IsCommandError(err error) bool {
	var tasErr *Error
	if errors.As(err, &tasErr) {
		return tasErr.Type == ErrorTypeCommand
	}
	return false
}

// IsParseError checks if the error is a parsing error.
func IsParseError(err error) bool {
	var tasErr *Error
	if errors.As(err, &tasErr) {
		return tasErr.Type == ErrorTypeParse
	}
	return false
}

// IsTimeoutError checks if the error is a timeout error.
func IsTimeoutError(err error) bool {
	var tasErr *Error
	if errors.As(err, &tasErr) {
		return tasErr.Type == ErrorTypeTimeout
	}
	return false
}

// IsDeviceError checks if the error is a device-specific error.
func IsDeviceError(err error) bool {
	var tasErr *Error
	if errors.As(err, &tasErr) {
		return tasErr.Type == ErrorTypeDevice
	}
	return false
}
