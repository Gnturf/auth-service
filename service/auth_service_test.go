package service

import (
	"auth-service/config"
	"auth-service/initializer"
	"auth-service/model/web"
	"auth-service/repository"
	"fmt"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

var configuration = &config.AppConfigInstance
var userRepository = repository.NewUserRepositoryImpl(configuration)
var tokenRepository = repository.NewTokenRepository(configuration)
var emailVerificationRepository = repository.NewVerificationRepositoryImpl()
var db = initializer.NewPostgressDB()
var redisConn = initializer.NewRedisConn(configuration)

func TestCheckEmailFound(t *testing.T) {

	authService := NewAuthService(userRepository, tokenRepository, emailVerificationRepository, configuration, db, redisConn)

	request := web.UserCheckEmailRequest{
		Email: "gfirmansyah279@gmail.com",
	}
	
	res, err := authService.CheckEmail(request)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, true, res)

	fmt.Println(res)
}

func TestCheckEmailNotFound(t *testing.T) {
	authService := NewAuthService(userRepository, tokenRepository, emailVerificationRepository, configuration, db, redisConn)

	request := web.UserCheckEmailRequest{
		Email: "gfirmansyah270@gmail.com",
	}
	
	res, err := authService.CheckEmail(request)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, false, res)

	fmt.Println(res)
}

func TestSignupNoMatchedEmail(t *testing.T) {
	authService := NewAuthService(userRepository, tokenRepository, emailVerificationRepository, configuration, db, redisConn)

	userCreate := web.UserCreateRequest{
		Email: "gfirmansyah270@gmail.com",
		Username: "JackieChan",
		Password: "JohnDoe^1280",
	}

	response, err := authService.Signup(userCreate)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
}