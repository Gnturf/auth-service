package controller

import (
	"auth-service/config"
	"auth-service/model/web"
	"auth-service/service"
	"fmt"

	"github.com/labstack/echo/v4"
)

/**
GENERAL STEP TO WRITE CONTROLLER
1. Bind the data to model
2. Basic validate the model
3. Pass the data to service
4. Create response object & set header
5. Send response object
*/

type AuthControllerImpl struct {
	AuthService service.AuthService
	Config *config.AppConfig
}

func NewAuthController(authService service.AuthService, config *config.AppConfig) AuthController {
	return &AuthControllerImpl{
		AuthService: authService,
		Config: config,
	}
}

// DONE 12/05/25
func (controller *AuthControllerImpl) HandleSignup(c echo.Context) error {
	// 1. Bind the email, password, profile name, image_url to UserCreateRequest struct
	var createRequest web.UserCreateRequest
	err := c.Bind(&createRequest)
	if err != nil {
		return c.JSON(200, web.WebResponse{
			Status: "failed",
			Code: 4000,
			Data: "Invalid request",
		})
	}

	// 2. Validate the request
	err = createRequest.BasicValidate()
	if err != nil {
		return c.JSON(200, web.WebResponse{
			Status: "failed",
			Code: 4000,
			Data: err.Error(),
		})
	}

	// 3. Pass UserCreateRequest to the authService to create a new user
	data, err := controller.AuthService.Signup(createRequest)
	if err != nil {
		return c.JSON(200, web.WebResponse{
			Status: "failed",
			Code: 4000,
			Data: err.Error(),
		})
	}

	// 4. Create response object with the user data and status code
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	webResponse := web.WebResponse{
		Code:  200,
		Status: "Success",
		Data:  data,
	}

	// 5. Return the response to the client
	return c.JSON(200, webResponse)
}

func (controller *AuthControllerImpl) HandleSignin(c echo.Context) error {
	// 1. Bind the email, password, profile name, image_url to UserSigninRequest struct
	userSigninRequest := web.UserSigninRequest{}
	err := c.Bind(&userSigninRequest)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	// 2. Validate the request
	err = userSigninRequest.BasicValidate()
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	// 3. Pass UserCreateRequest to the authService to create a new user
	data, err := controller.AuthService.Signin(userSigninRequest)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	// 4. Create response object with the user data and status code
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	webResponse := web.WebResponse{
		Code:  200,
		Status: "Success",
		Data:  data,
	}

	return c.JSON(200, webResponse)
}

// DONE 12/05/25
func (controller *AuthControllerImpl) HandleCheckEmail(c echo.Context) error {
	// 1. Bind the data to model
	userCheckEmailRequest := web.UserCheckEmailRequest{}
	err := c.Bind(&userCheckEmailRequest)
	if err != nil {
		return c.JSON(200, web.WebResponse{Status: "failed", Code: 4000, Data: "Invalid request"})
	}

	// 2. Basic validate the model
	err = userCheckEmailRequest.BasicValidate()
	if err != nil {
		return c.JSON(200, web.WebResponse{Status: "failed", Code: 4000, Data: err.Error()})
	}

	// 3. Pass the data to service
	res, err := controller.AuthService.CheckEmail(userCheckEmailRequest)
	if err != nil {
		return c.JSON(200, web.WebResponse{Status: "failed", Code: 4000, Data: err.Error()})
	}

	// 4. Create response object & set header
	var response web.WebResponse

	if res {
		response = web.WebResponse{
			Status: "failed",
			Code: 4009,
			Data: "email already exist",
		}
	} else {
		response = web.WebResponse{
			Status: "success",
			Code: 200,
		}
	}

	// 5. Send response object
	return c.JSON(200, response)
}

func (controller *AuthControllerImpl) HandleRefresh(c echo.Context) error {
	// 1. Bind the data
	userRefreshToken := web.UserRefreshRequest{}
	err := c.Bind(&userRefreshToken)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	fmt.Println(userRefreshToken)

	// 2. Validate the request
	err = userRefreshToken.BasicValidate()
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	// 3. Pass userTokenRefresh to the authService 
	data, err := controller.AuthService.Refresh(userRefreshToken)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	// 4. Create response object and status code
	response := web.WebResponse{
		Code: 200,
		Status: "Success",
		Data: data,
	}
	
	// 5. Return the response to the client
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return c.JSON(200, response)
}

func (controller *AuthControllerImpl) HandleSignout(c echo.Context) error {
	// 1. Bind the data to models
	signOutRequest := web.UserSignoutRequest{}
	if err := c.Bind(&signOutRequest); err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	// 2. Perform basic validation
	if err := signOutRequest.BasicValidate(); err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	// 3. Pass data to service
	if err := controller.AuthService.Signout(signOutRequest); err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	// 4. Create response struct and set header
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := web.WebResponse{
		Code: 200,
		Status: "success",
		Data: nil,
	}

	// 5. Send back to user
	return c.JSON(200, response)
}