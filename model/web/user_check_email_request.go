package web

import (
	"auth-service/exception"
	"auth-service/helper"
)

type UserCheckEmailRequest struct {
	Email string `form:"email" json:"email" `
}

func (s *UserCheckEmailRequest) BasicValidate() error {
	// Validate if the fields are empty
	if err := helper.ValidateIsEmpty(s.Email, "Email"); err != nil {
		return helper.ErrorHandler(err, exception.EmailCannotBeEmpty)
	}

	return nil
}

func (s *UserCheckEmailRequest) Validate() error {
	// Validate if the fields are empty
	if err := helper.ValidateIsEmpty(s.Email, "Email"); err != nil {
		return helper.ErrorHandler(err, exception.EmailCannotBeEmpty)
	}

	// Validate fields structures
	if err := helper.ValidateEmail(s.Email); err != nil {
		return helper.ErrorHandler(err, exception.InvalidEmail)
	}

	return nil
}