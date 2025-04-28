package helper

import (
	"errors"
)

func ErrorHandler(err error, errorMessage string) error {
	if err != nil {
		return errors.New(errorMessage)
	}

	return nil
}