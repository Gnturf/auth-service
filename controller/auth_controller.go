package controller

import "github.com/labstack/echo/v4"

type AuthController interface {
	HandleSignup(c echo.Context) error
	HandleSignin(c echo.Context) error
	HandleRefresh(c echo.Context) error
	HandleSignout(c echo.Context) error
}