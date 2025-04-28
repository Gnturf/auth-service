package web

import (
	"auth-service/exception"
	"auth-service/helper"
)

type UserSigninRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (s *UserSigninRequest) BasicValidate() error {
	// Validate if the fields are empty
	if err := helper.ValidateIsEmpty(s.Email, "Email"); err != nil {
		return helper.ErrorHandler(err, exception.EmailCannotBeEmpty)
	}

	if err := helper.ValidateIsEmpty(s.Password, "Password"); err != nil {
		return helper.ErrorHandler(err, exception.PasswordCannotBeEmpty)
	}

	return nil
}

func (s *UserSigninRequest) Validate() error {
	// Validate if the fields are empty
	if err := helper.ValidateIsEmpty(s.Email, "Email"); err != nil {
		return helper.ErrorHandler(err, exception.EmailCannotBeEmpty)
	}

	if err := helper.ValidateIsEmpty(s.Password, "Password"); err != nil {
		return helper.ErrorHandler(err, exception.PasswordCannotBeEmpty)
	}

	// Validate fields structures
	if err := helper.ValidateEmail(s.Email); err != nil {
		return helper.ErrorHandler(err, exception.InvalidEmail)
	}

	if err := helper.ValidatePassword(s.Password); err != nil {
		return helper.ErrorHandler(err, err.Error())
	}

	return nil
}