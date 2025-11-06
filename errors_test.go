package tasmota

import (
	"errors"
	"testing"
)

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected string
	}{
		{"network", ErrorTypeNetwork, "network"},
		{"auth", ErrorTypeAuth, "auth"},
		{"command", ErrorTypeCommand, "command"},
		{"parse", ErrorTypeParse, "parse"},
		{"timeout", ErrorTypeTimeout, "timeout"},
		{"device", ErrorTypeDevice, "device"},
		{"unknown", ErrorType(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.errType.String(); got != tt.expected {
				t.Errorf("ErrorType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name: "with wrapped error",
			err: &Error{
				Type:    ErrorTypeNetwork,
				Message: "connection failed",
				Err:     errors.New("connection refused"),
			},
			expected: "network error: connection failed: connection refused",
		},
		{
			name: "without wrapped error",
			err: &Error{
				Type:    ErrorTypeAuth,
				Message: "invalid credentials",
				Err:     nil,
			},
			expected: "auth error: invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("underlying error")
	err := &Error{
		Type:    ErrorTypeNetwork,
		Message: "test",
		Err:     wrappedErr,
	}

	if got := err.Unwrap(); got != wrappedErr {
		t.Errorf("Error.Unwrap() = %v, want %v", got, wrappedErr)
	}
}

func TestNewError(t *testing.T) {
	wrappedErr := errors.New("wrapped")
	err := NewError(ErrorTypeCommand, "test message", wrappedErr)

	if err.Type != ErrorTypeCommand {
		t.Errorf("NewError() Type = %v, want %v", err.Type, ErrorTypeCommand)
	}
	if err.Message != "test message" {
		t.Errorf("NewError() Message = %v, want %v", err.Message, "test message")
	}
	if err.Err != wrappedErr {
		t.Errorf("NewError() Err = %v, want %v", err.Err, wrappedErr)
	}
}

func TestIsNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "network error",
			err:      NewError(ErrorTypeNetwork, "test", nil),
			expected: true,
		},
		{
			name:     "auth error",
			err:      NewError(ErrorTypeAuth, "test", nil),
			expected: false,
		},
		{
			name:     "standard error",
			err:      errors.New("test"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNetworkError(tt.err); got != tt.expected {
				t.Errorf("IsNetworkError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsAuthError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "auth error",
			err:      NewError(ErrorTypeAuth, "test", nil),
			expected: true,
		},
		{
			name:     "network error",
			err:      NewError(ErrorTypeNetwork, "test", nil),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAuthError(tt.err); got != tt.expected {
				t.Errorf("IsAuthError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsCommandError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "command error",
			err:      NewError(ErrorTypeCommand, "test", nil),
			expected: true,
		},
		{
			name:     "network error",
			err:      NewError(ErrorTypeNetwork, "test", nil),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCommandError(tt.err); got != tt.expected {
				t.Errorf("IsCommandError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsParseError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "parse error",
			err:      NewError(ErrorTypeParse, "test", nil),
			expected: true,
		},
		{
			name:     "network error",
			err:      NewError(ErrorTypeNetwork, "test", nil),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsParseError(tt.err); got != tt.expected {
				t.Errorf("IsParseError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsTimeoutError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "timeout error",
			err:      NewError(ErrorTypeTimeout, "test", nil),
			expected: true,
		},
		{
			name:     "network error",
			err:      NewError(ErrorTypeNetwork, "test", nil),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTimeoutError(tt.err); got != tt.expected {
				t.Errorf("IsTimeoutError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsDeviceError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "device error",
			err:      NewError(ErrorTypeDevice, "test", nil),
			expected: true,
		},
		{
			name:     "network error",
			err:      NewError(ErrorTypeNetwork, "test", nil),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDeviceError(tt.err); got != tt.expected {
				t.Errorf("IsDeviceError() = %v, want %v", got, tt.expected)
			}
		})
	}
}
