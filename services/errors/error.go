package errors

import (
	"fmt"
)

// Define custom error types
type InternalError struct {
	Message string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Internal error: %s", e.Message)
}

type KeyNotFoundError struct {
	Key string
}

func (e *KeyNotFoundError) Error() string {
	return fmt.Sprintf("Ticker not found: %s", e.Key)
}
