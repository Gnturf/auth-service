package exception

import (
	"errors"
	"fmt"
)

const (
	UserAlreadyExists							= "user already exist"
	UserWasNotFound								= "user was not found"
	EmailWasNotFound							= "email was not found"
	EmailCannotBeEmpty        		= "email cannot be empty"
	PasswordCannotBeEmpty     		= "password cannot be empty"
	UsernameCannotBeEmpty     		= "username cannot be empty"
	InvalidEmail              		= "email is not valid"
	InvalidPassword      		   		= "password is not valid"
	InvalidUsername           		= "username is not valid"
	EmailAlreadyExists       			= "email already exists"
	UsernameAlreadyExists     		= "username already exists"
	UsernameWasNotFound						= "username was not found"
	UsernameViolateRegexRules 		= "username can only contain letters, numbers, and underscores"
	RefreshTokenCannotBeEmpty 		= "refresh token cannot be empty"
	RefreshTokenViolateRegexRules = "refresh token is not valid"
	RefreshTokenWasNotFound 			= "refresh token was not found"
	AccessTokenFailedToGenerated 	= "failed to generate access token"
)

var (
	ErrUserWasNotFound = errors.New(UserWasNotFound)
	ErrUserAlreadyExist = errors.New(UserAlreadyExists)
	ErrEmailWasNotFound = errors.New(EmailWasNotFound)
	ErrUsernameWasNotFound = errors.New(UsernameWasNotFound)
	ErrEmailAlreadyExist = errors.New(EmailAlreadyExists)
	ErrRefreshTokenWasNotFound = errors.New(RefreshTokenWasNotFound)
)

func FieldMinMaxCharacters(field string, min int, max int) string {
	return field + " must be between " + fmt.Sprint(min) + " and " + fmt.Sprint(max) + " characters"
}

func FieldMinCharacter(field string, min int) string {
	return field + " must be at least " + fmt.Sprint(min) + " characters long"
}