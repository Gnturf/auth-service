package main

import (
	"auth-service/config"
	"auth-service/controller"
	"auth-service/initializer"
	"auth-service/repository"
	"auth-service/service"
	"log"
	"os"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

/*
	AUTH SERVICE
	1. /auth/signin
	2. /auth/signup
	3. /auth/refresh-token
	4. /auth/logout
	5. /auth/verify-email
*/

func init() {
	err := godotenv.Load()
  if err != nil {
      log.Fatal("Error loading .env file")
  }

	config.LoadEnv()
}

func main() {
	e := echo.New()

	logger := log.Default()
	db := initializer.NewPostgressDB()
	validate := validator.New()

	authRepository := repository.NewAuthRepository(&config.AppConfigInstance, logger)
	tokenRepository := repository.NewTokenRepository(&config.AppConfigInstance, logger)
	service := service.NewAuthService(authRepository, tokenRepository, &config.AppConfigInstance, logger, db, validate)
	controller := controller.NewAuthController(service, &config.AppConfigInstance)

	authGroup := e.Group("/auth")
	authGroup.POST("/signup", controller.HandleSignup)
	authGroup.POST("/signin", controller.HandleSignin)
	authGroup.POST("/refresh", controller.HandleRefresh)
	authGroup.POST("/signout", controller.HandleSignout)

	port := os.Getenv("PORT")
	if port == "" {
			log.Fatal("PORT environment variable not set")
	}

	logger.Printf("Starting server")

	e.Logger.Fatal(e.Start(":" +port))
}