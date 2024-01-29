package errors

import (
	"errors"
	"fmt"
)

const (
	ListingExist        = "Symbol %s exists."
	ListingDoesNotExist = "Symbol %s exists. Put symbol to continue."
	FailedToReadData    = "Failed to read data from database: %s."
)

func GetErrorMessage(template string, values ...string) error {
	errMessage := fmt.Sprintf(template, values)
	return errors.New(errMessage)
}
