package repository

import (
	"auth-service/config"
	"auth-service/exception"
	"auth-service/initializer"
	"auth-service/model/domain"
	"context"
	"errors"
	"fmt"
	"testing"
)

// func TestErrorComparasion(t *testing.T) {
// 	// var notFound = errors.New(exception.UserWasNotFound)
// 	// Create Repository Error
// 	err := exception.NewRepositoryError(exception.UserWasNotFound, errors.New(exception.UserWasNotFound) , "Fetch")

// 	if errors.Is(err, errors.New(exception.UserWasNotFound)) {
// 		fmt.Println("Same")
// 		return
// 	}

// 	fmt.Println("Not the same")
// }

func TestFetchUserByEmailNotFound(t *testing.T) {
	ctx := context.Background()
	userDomain := domain.User{
		Email: "aaaaaaaaaa@gmail.com",
	}

	db := initializer.NewPostgressDB()
	tx, txErr := db.Begin()
	if txErr != nil {
		panic(txErr)
	}

	userRepository := NewUserRepositoryImpl(&config.AppConfigInstance)

	_, err := userRepository.FetchUserByEmail(ctx, &userDomain, tx)
	if err != nil {
		txErr = err
		if errors.Is(txErr, exception.ErrUserWasNotFound) {
			fmt.Println("User Was not Found")
			return 
		}

		panic(err)
	}

	fmt.Println("User Found")
}

func TestFetchUserByEmailFound(t *testing.T) {
	ctx := context.Background()
	userDomain := domain.User{
		Email: "2KaFQzaND2@gmail.com",
	}

	db := initializer.NewPostgressDB()
	tx, txErr := db.Begin()
	if txErr != nil {
		panic(txErr)
	}

	userRepository := NewUserRepositoryImpl(&config.AppConfigInstance)

	_, err := userRepository.FetchUserByEmail(ctx, &userDomain, tx)
	if err != nil {
		txErr = err
		if errors.Is(txErr, exception.ErrUserWasNotFound) {
			fmt.Println("User Was not Found")
			return 
		}

		panic(err)
	}

	fmt.Println("User Found")
}