package web

import (
	"auth-service/exception"
	"auth-service/helper"
)

type UserCreateRequest struct {
	Email    string `form:"email" json:"email" `
	Password string `form:"password" json:"password" `
	Username string `form:"username" json:"username" `
}

func (s *UserCreateRequest) BasicValidate() error {
	// Validate if the fields are empty
	if err := helper.ValidateIsEmpty(s.Email, "Email"); err != nil {
		return helper.ErrorHandler(err, exception.EmailCannotBeEmpty)
	}

	if err := helper.ValidateIsEmpty(s.Password, "Password"); err != nil {
		return helper.ErrorHandler(err, exception.PasswordCannotBeEmpty)
	}

	if err := helper.ValidateIsEmpty(s.Username, "Username"); err != nil {
		return helper.ErrorHandler(err, exception.UsernameCannotBeEmpty)
	}

	return nil
}

func (s *UserCreateRequest) Validate() error {

	// Validate if the fields are empty
	if err := helper.ValidateIsEmpty(s.Email, "Email"); err != nil {
		return helper.ErrorHandler(err, exception.EmailCannotBeEmpty)
	}

	if err := helper.ValidateIsEmpty(s.Password, "Password"); err != nil {
		return helper.ErrorHandler(err, exception.PasswordCannotBeEmpty)
	}

	if err := helper.ValidateIsEmpty(s.Username, "Username"); err != nil {
		return helper.ErrorHandler(err, exception.UsernameCannotBeEmpty)
	}


	// Validate fields structures
	if err := helper.ValidateUsername(s.Username); err != nil {
		return helper.ErrorHandler(err, err.Error())
	}

	if err := helper.ValidateEmail(s.Email); err != nil {
		return helper.ErrorHandler(err, exception.InvalidEmail)
	}

	if err := helper.ValidatePassword(s.Password); err != nil {
		return helper.ErrorHandler(err, err.Error())
	}

	return nil
}