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

// Function to create custom errors
func NewInternalError(message string) error {
	return &InternalError{Message: message}
}

func NewKeyNotFoundError(key string) error {
	return &KeyNotFoundError{Key: key}
}

// Function to check error types
func CheckError(err error) {
	switch err.(type) {
	case *InternalError:
		fmt.Println("Internal error occurred:", err.Error())
	case *KeyNotFoundError:
		fmt.Println("Key not found:", err.Error())
	default:
		fmt.Println("Unknown error:", err.Error())
	}
}
