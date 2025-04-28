package service

import "auth-service/model/web"

type AuthService interface {
	Signup(user web.UserCreateRequest) (web.UserCreateResponse, error)
	Signin(user web.UserSigninRequest) (web.UserSigninResponse, error)
	Refresh(user web.UserRefreshRequest) (web.UserRefreshResponse, error)
	Signout(user web.UserSignoutRequest) (error)
}