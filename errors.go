package gigachat

import "fmt"

type GigaChatError struct {
	Message string
	Code    int
	Err     error
}

func (e *GigaChatError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("gigachat error (code %d): %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("gigachat error (code %d): %s", e.Code, e.Message)
}

func (e *GigaChatError) Unwrap() error {
	return e.Err
}

type AuthenticationError struct {
	Message string
	Err     error
}

func (e *AuthenticationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("authentication error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("authentication error: %s", e.Message)
}

func (e *AuthenticationError) Unwrap() error {
	return e.Err
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}
