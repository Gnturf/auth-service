package main

import (
	"auth-service/config"
	"auth-service/controller"
	"auth-service/initializer"
	"auth-service/repository"
	"auth-service/service"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

/*
	AUTH SERVICE
	1. /check-email ✔️
	2. /signup ✔️
	3. /signin
	4. /signout
	5. /refresh-token
	6. /verify-email
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
	redis := initializer.NewRedisConn(&config.AppConfigInstance)

	userRepository := repository.NewUserRepositoryImpl(&config.AppConfigInstance)
	tokenRepository := repository.NewTokenRepository(&config.AppConfigInstance)
	emailVerificationRepository := repository.NewVerificationRepositoryImpl()
	service := service.NewAuthService(userRepository, tokenRepository, emailVerificationRepository, &config.AppConfigInstance, db, redis)
	controller := controller.NewAuthController(service, &config.AppConfigInstance)

	e.POST("/signup", controller.HandleSignup)
	e.POST("/signin", controller.HandleSignin)
	e.POST("/check-email", controller.HandleCheckEmail)
	e.POST("/refresh", controller.HandleRefresh)
	e.POST("/signout", controller.HandleSignout)

	port := os.Getenv("PORT")
	if port == "" {
			log.Fatal("PORT environment variable not set")
	}

	logger.Printf("Starting server")

	e.Logger.Fatal(e.Start(":" +port))
}