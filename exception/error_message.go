package exception

import "fmt"

const (
	UserAlreadyExists							= "user already exist"
	UserWasNotFound								= "user was not found"
	EmailCannotBeEmpty        		= "email cannot be empty"
	PasswordCannotBeEmpty     		= "password cannot be empty"
	UsernameCannotBeEmpty     		= "username cannot be empty"
	InvalidEmail              		= "email is not valid"
	InvalidPassword      		   		= "password is not valid"
	InvalidUsername           		= "username is not valid"
	EmailAlreadyExists       			= "email already exists"
	UsernameAlreadyExists     		= "username already exists"
	UsernameViolateRegexRules 		= "username can only contain letters, numbers, and underscores"
	RefreshTokenCannotBeEmpty 		= "refresh token cannot be empty"
	RefreshTokenViolateRegexRules = "refresh token is not valid"
	RefreshTokenWasNotFound 			= "refresh token was not found"
	AccessTokenFailedToGenerated 	= "failed to generate access token"
)

func FieldMinMaxCharacters(field string, min int, max int) string {
	return field + " must be between " + fmt.Sprint(min) + " and " + fmt.Sprint(max) + " characters"
}

func FieldMinCharacter(field string, min int) string {
	return field + " must be at least " + fmt.Sprint(min) + " characters long"
}