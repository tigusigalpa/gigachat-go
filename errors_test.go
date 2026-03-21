package gigachat

import (
	"errors"
	"testing"
)

func TestGigaChatError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *GigaChatError
		expected string
	}{
		{
			name: "error with code and wrapped error",
			err: &GigaChatError{
				Message: "request failed",
				Code:    500,
				Err:     errors.New("connection timeout"),
			},
			expected: "gigachat error (code 500): request failed: connection timeout",
		},
		{
			name: "error with code only",
			err: &GigaChatError{
				Message: "invalid request",
				Code:    400,
			},
			expected: "gigachat error (code 400): invalid request",
		},
		{
			name: "error with zero code",
			err: &GigaChatError{
				Message: "unknown error",
				Code:    0,
			},
			expected: "gigachat error (code 0): unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("GigaChatError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGigaChatError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := &GigaChatError{
		Message: "test",
		Code:    500,
		Err:     wrappedErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != wrappedErr {
		t.Errorf("GigaChatError.Unwrap() = %v, want %v", unwrapped, wrappedErr)
	}

	errNoWrap := &GigaChatError{
		Message: "test",
		Code:    500,
	}

	if unwrapped := errNoWrap.Unwrap(); unwrapped != nil {
		t.Errorf("GigaChatError.Unwrap() = %v, want nil", unwrapped)
	}
}

func TestAuthenticationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AuthenticationError
		expected string
	}{
		{
			name: "error with wrapped error",
			err: &AuthenticationError{
				Message: "invalid credentials",
				Err:     errors.New("401 unauthorized"),
			},
			expected: "authentication error: invalid credentials: 401 unauthorized",
		},
		{
			name: "error without wrapped error",
			err: &AuthenticationError{
				Message: "token expired",
			},
			expected: "authentication error: token expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("AuthenticationError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAuthenticationError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := &AuthenticationError{
		Message: "test",
		Err:     wrappedErr,
	}

	if unwrapped := err.Unwrap(); unwrapped != wrappedErr {
		t.Errorf("AuthenticationError.Unwrap() = %v, want %v", unwrapped, wrappedErr)
	}

	errNoWrap := &AuthenticationError{
		Message: "test",
	}

	if unwrapped := errNoWrap.Unwrap(); unwrapped != nil {
		t.Errorf("AuthenticationError.Unwrap() = %v, want nil", unwrapped)
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "validation error with message",
			err: &ValidationError{
				Message: "field is required",
			},
			expected: "validation error: field is required",
		},
		{
			name: "validation error with empty message",
			err: &ValidationError{
				Message: "",
			},
			expected: "validation error: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("ValidationError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}
