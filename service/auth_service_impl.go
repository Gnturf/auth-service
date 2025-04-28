package service

import (
	"auth-service/config"
	"auth-service/exception"
	"auth-service/helper"
	"auth-service/model/domain"
	"auth-service/model/web"
	"auth-service/repository"
	"context"
	"database/sql"
	"log"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	AuthRepository repository.AuthRepository
	TokenRepository repository.TokenRepository
	Config *config.AppConfig
	Logger *log.Logger
	DB *sql.DB
}

func NewAuthService(authRepository repository.AuthRepository, tokenRepository repository.TokenRepository, config *config.AppConfig, logger *log.Logger, db *sql.DB, validation *validator.Validate) AuthService {
	return &AuthServiceImpl{
		AuthRepository: authRepository,
		TokenRepository: tokenRepository,
		Config: config,
		Logger: logger,
		DB: db,
	}
}

func (service *AuthServiceImpl) Signup(user web.UserCreateRequest) (web.UserCreateResponse, error) {
	// Creating empty response
	response := web.UserCreateResponse{}
	service.Logger.Printf("Creating empty UserCreateResponse")

	// Mapping data from Request to Domain
	service.Logger.Printf("Map data to user domain")	
	userDomain := domain.User{
		Email:     	user.Email,
		Password: 	user.Password,
		Username:  	user.Username,
	}
	service.Logger.Printf("Map with Email: %s, Password %s, Username: %s", userDomain.Email, userDomain.Password, userDomain.Username)

	ctx := context.Background()
	
	var txErr error

	// Begin DB Connection
	service.Logger.Printf("[%s] Begin DB Connection", user.Email)
	tx, txErr := service.DB.Begin()
	if txErr != nil {
		service.Logger.Printf("[%s] Failed Begin DB Connection", user.Email)
		return response, exception.NewServiceError(txErr.Error(), "Begin DB Conn")
	}

	// Create Defer
	defer func() {
		if txErr != nil {
			service.Logger.Printf("[%s] Transaction Rollback", user.Email)
			tx.Rollback()	
		} else if recover() != nil {
			service.Logger.Printf("[%s] Transaction Rollback", user.Email)
			tx.Rollback()
		} else { 
			service.Logger.Printf("[%s] Transaction Commited", user.Email)
			tx.Commit()
		}
	}()

	// 1.B. If user does not exist create a new user in auth-service
	// 1.B.1. Fully validate the request
	txErr = user.Validate()
	if txErr != nil {
		txErr = exception.NewServiceError(txErr.Error(), "User Validation")
		return response, txErr
	}
	
	// 1.B.2. Hash the password using bcrypt
	var hashPassword []byte
	hashPassword, txErr = bcrypt.GenerateFromPassword([]byte(user.Password), service.Config.BcryptHashCost)
	if txErr != nil {
		txErr = exception.NewServiceError(txErr.Error(), "Password Hashing")
		return response, txErr
	}

	userDomain.Password = string(hashPassword)

	// 1.B.3. Create a new user in the database
	var uuid string
	uuid, txErr = service.AuthRepository.CreateUser(ctx, userDomain, tx)
	if txErr != nil {
		return response, txErr
	}

	userDomain.Id = uuid

	// 2. Generate a JWT access token
	signedToken, txErr := helper.GenerateJWTAccessToken(uuid, service.Config.SecretKey)
	if txErr != nil {
		txErr = exception.NewServiceError(txErr.Error(), "Token Signing")
		return response, txErr
	}

	// 3. Generate refresh 
	var refreshToken string
	refreshToken, txErr = helper.GenerateRefreshToken(service.Config.RefreshTokenLength)
	if txErr != nil {
		txErr = exception.NewServiceError(txErr.Error(), "Refresh Token Generation")
		return response, txErr
	}

	// 4. Store the refresh token in the database
	txErr = service.TokenRepository.InsertRefreshToken(ctx, userDomain, refreshToken, tx)
	if txErr != nil {
		return response, txErr
	}

	// 5. Send the access token and refresh token in the response
	response.AccessToken = signedToken
	response.RefreshToken = refreshToken
	response.UUID = uuid

	return response, nil
}

func (service *AuthServiceImpl) Signin(user web.UserSigninRequest) (web.UserSigninResponse, error) {
	response := web.UserSigninResponse{}
	userDomain := domain.User{
		Email:    	user.Email,
		Password: 	user.Password,
	}

	ctx := context.Background()
	var txErr error

	tx, txErr := service.DB.Begin()
	if txErr != nil {
		return response, exception.NewServiceError(txErr.Error(), "Transaction Begin")
	}

	defer func() {
		if txErr != nil {
			tx.Rollback()
		} else if recover() != nil {
			tx.Rollback()
		} else { 
			tx.Commit()
		}
	}()

	// 1. Fully validate the request
	txErr = user.Validate()
	if txErr != nil {
		txErr = exception.NewServiceError(txErr.Error(), "User Validation")
		return response, txErr
	}

	// 2. Check if the user exists in the database
	_, err := service.AuthRepository.FetchUser(ctx, &userDomain, tx)
	if err != nil {
		txErr = err
		return response, err
	}

	// 3. Compare the password with the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(userDomain.Password), []byte(user.Password))
	if err != nil {
		txErr = exception.NewServiceError(exception.InvalidPassword, "Password Comparison")
		return response, txErr
	}

	// 4. Generate a JWT access token 
	signedToken, txErr := helper.GenerateJWTAccessToken(userDomain.Id, service.Config.SecretKey)
	if txErr != nil {
		txErr = exception.NewServiceError(txErr.Error(), "Token Signing")
		return response, txErr
	}

	// 5. Generate refresh 
	var refreshToken string
	refreshToken, txErr = helper.GenerateRefreshToken(service.Config.RefreshTokenLength)
	if txErr != nil {
		txErr = exception.NewServiceError(txErr.Error(), "Refresh Token Generation")
		return response, txErr
	}

	// 6. Store the refresh token in the database
	txErr = service.TokenRepository.InsertRefreshToken(ctx, userDomain, refreshToken, tx)
	if txErr != nil {
		return response, txErr
	}

	// 7. Send the access token and refresh token in the response
	response.AccessToken = signedToken
	response.RefreshToken = refreshToken
	response.UUID = userDomain.Id

	return response, nil
}

func (service *AuthServiceImpl) Refresh(user web.UserRefreshRequest) (web.UserRefreshResponse, error) {
	response := web.UserRefreshResponse{}
	tokenDomain := domain.Tokens{
		RefreshToken: user.RefreshToken,
	}
	ctx := context.Background()
	var txErr error

	tx, txErr := service.DB.Begin()
	if txErr != nil {
		return response, exception.NewServiceError(txErr.Error(), "Transaction Begin")
	}

	defer func() {
		if txErr != nil {
			tx.Rollback()
		} else if recover() != nil {
			tx.Rollback()
		} else { 
			tx.Commit()
		}
	}()

	// Check token if exist and expired or not
	_, err := service.TokenRepository.FetchRefreshToken(ctx, &tokenDomain, tx)
	if err != nil {
		txErr = err
		return response, txErr
	}

	// Generate new access token
	var newAccessToken string
	newAccessToken, txErr = helper.GenerateJWTAccessToken(tokenDomain.Id, service.Config.SecretKey)
	if txErr != nil {
		return response, exception.NewServiceError(exception.AccessTokenFailedToGenerated, "Generate Access Token")
	}

	// Bind the tokenDomain to response
	response.UUID = tokenDomain.Id
	response.AccessToken = newAccessToken
	response.RefreshToken = tokenDomain.RefreshToken

	return response, nil
}

func (service *AuthServiceImpl) Signout(user web.UserSignoutRequest) (error) {
	// Mapping data from Request to Domain
	service.Logger.Printf("Map data to user domain")	
	tokenDomain := domain.Tokens{
		RefreshToken: user.RefreshToken,
	}
	service.Logger.Printf("Map with Refresh Token: %s", tokenDomain.RefreshToken)

	ctx := context.Background()
	var txErr error

	// Begin DB Connection
	service.Logger.Printf("[%s] Begin DB Connection", user.RefreshToken)
	tx, txErr := service.DB.Begin()
	if txErr != nil {
		service.Logger.Printf("[%s] Failed Begin DB Connection", user.RefreshToken)
		return exception.NewServiceError(txErr.Error(), "Begin DB Conn")
	}

	// Create Defer
	defer func() {
		if txErr != nil {
			service.Logger.Printf("[%s] Transaction Rollback", user.RefreshToken)
			tx.Rollback()	
		} else if recover() != nil {
			service.Logger.Printf("[%s] Transaction Rollback", user.RefreshToken)
			tx.Rollback()
		} else { 
			service.Logger.Printf("[%s] Transaction Commited", user.RefreshToken)
			tx.Commit()
		}
	}()

	err := service.TokenRepository.RevokeRefreshToken(ctx, &tokenDomain, tx);
	if err != nil {
		txErr = err
		return txErr
	}

	return nil
}