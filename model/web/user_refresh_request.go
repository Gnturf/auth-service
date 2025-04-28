package web

import (
	"auth-service/exception"
	"auth-service/helper"
)

type UserRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (s *UserRefreshRequest) BasicValidate() error {
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