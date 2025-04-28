package web

import (
	"auth-service/exception"
	"auth-service/helper"
)

type UserSignoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *UserSignoutRequest) BasicValidate() error {
	// Validate if the fields are empty
	if err := helper.ValidateIsEmpty(s.RefreshToken, "RefreshToken"); err != nil {
		return helper.ErrorHandler(err, exception.RefreshTokenCannotBeEmpty)
	}

	// Validate fields structures
	if err := helper.ValiadteAlphanumericOnly(s.RefreshToken); err != nil {
		return helper.ErrorHandler(err, exception.RefreshTokenViolateRegexRules)
	}

	return nil
}