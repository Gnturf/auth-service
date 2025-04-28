package helper

import (
	"auth-service/exception"
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateIsEmpty(value string, fieldName string) error {
	if value == "" {
		return errors.New("")
	}
	return nil
}

func ValidateUsername(username string) error {
	minLength := 8

	if len(username) < minLength {
		return errors.New(exception.FieldMinCharacter("Username", minLength))
	}

	// Use ozzo-validation to apply multiple rules for the username.
	err := validation.Validate(username,
		// Ensure the username is not empty
		validation.Required.Error(exception.UsernameCannotBeEmpty),

		// Ensure the username matches the pattern (alphanumeric + underscore)
		validation.Match(regexp.MustCompile(`^[a-zA-Z0-9_]+$`)).Error(exception.UsernameViolateRegexRules),
	)

	if err != nil {
		// If validation fails, return an error
		return err
	}

	return nil
}

func ValidateEmail(email string) error {
	err := validation.Validate(email, is.Email)
	if err != nil {
		return err
	}

	return nil
}

func ValidatePassword(password string) error {
	var (
		minLength        = 8
		uppercasePattern = `[A-Z]`
		lowercasePattern = `[a-z]`
		numberPattern    = `[0-9]`
		specialPattern   = `[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`
	)

	if len(password) < minLength {
		return errors.New("password must be at least 8 characters long")
	}

	if match, _ := regexp.MatchString(uppercasePattern, password); !match {
		return errors.New("password must contain at least one uppercase letter")
	}

	if match, _ := regexp.MatchString(lowercasePattern, password); !match {
		return errors.New("password must contain at least one lowercase letter")
	}

	if match, _ := regexp.MatchString(numberPattern, password); !match {
		return errors.New("password must contain at least one number")
	}

	if match, _ := regexp.MatchString(specialPattern, password); !match {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func ValiadteAlphanumericOnly(value string) error {
	// Define a regular expression pattern for alphanumeric characters
	pattern := `^[a-zA-Z0-9]+$`

	// Check if the value matches the pattern
	if match, _ := regexp.MatchString(pattern, value); !match {
		return errors.New("value can only contain letters and numbers")
	}

	return nil
}