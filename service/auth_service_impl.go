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
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

/**
GENERAL STEP TO WRITE SERVICE
1. Create necessary variable (e.g context, db)
2. Create defer for db
3. Fully validate the request
4. Bind the data to domain
5. custom ...
6. Create response
*/

type AuthServiceImpl struct {
	UserRepository repository.UserRepository
	TokenRepository repository.TokenRepository
	EmailVerificationRepository repository.EmailVerificationRepository
	Config *config.AppConfig
	DB *sql.DB
	Redis *redis.Client
}

func NewAuthService(userRepository repository.UserRepository, tokenRepository repository.TokenRepository, emailVerificationRepository repository.EmailVerificationRepository, config *config.AppConfig, db *sql.DB, redis *redis.Client) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		TokenRepository: tokenRepository,
		EmailVerificationRepository: emailVerificationRepository,
		Config: config,
		DB: db,
		Redis: redis,
	}
}

// DONE 12/05/25
func (service *AuthServiceImpl) Signup(request web.UserCreateRequest) (web.UserCreateResponse, error) {
	// 1. Create necessary variable (e.g context, db)
	response := web.UserCreateResponse{}
	var txErr error
	ctx := context.Background()

	tx, txErr := service.DB.Begin()
	if txErr != nil {
		return response, exception.NewServiceError(txErr.Error(), "Transaction Begin")
	}

	// 2. Create defer for db
	defer func() {
		if txErr != nil {
			tx.Rollback()	
		} else if recover() != nil {
			tx.Rollback()
		} else { 
			tx.Commit()
		}
	}()

	// 3. Fully validate the request
	txErr = request.Validate()
	if txErr != nil {
		txErr = exception.NewServiceError(txErr.Error(), "User Validation")
		return response, txErr
	}

	// 4. Bind the data to domain
	userDomain := domain.User{
		Email:     	request.Email,
		Password: 	request.Password,
		Username:  	request.Username,
	}


	// 5. custom ...
	// 5.1 Check if user with email already exist
	_, err := service.UserRepository.FetchUserByEmail(ctx, &userDomain, tx)
	if err != nil {
		txErr = err
		if errors.Is(txErr, exception.ErrEmailWasNotFound) {
			
		} else {
			return response, txErr
		}
	} else {
		return response, exception.NewServiceError(exception.EmailAlreadyExists, "Email Check")
	}

	// 5.2 Check if user with username already exist
	_, err = service.UserRepository.FetchUserByUsername(ctx, &userDomain, tx)
	if err != nil {
		txErr = err
		if errors.Is(txErr, exception.ErrUsernameWasNotFound) {
			
		} else {
			return response, txErr
		}
	} else {
		return response, exception.NewServiceError(exception.UsernameAlreadyExists, "Email Check")
	}

	// 5.3 Generate Email Verification Token
	emailVerificationToken, err := helper.GenerateNumericToken(6)
	if err != nil {
		return response, exception.NewServiceError(err.Error(), "Email Check")
	}

	// 5.4 Store user db:users email_verified:false
	// 5.4.A Hash the password using bcrypt
	var hashPassword []byte
	hashPassword, txErr = bcrypt.GenerateFromPassword([]byte(request.Password), service.Config.BcryptHashCost)
	if txErr != nil {
		txErr = exception.NewServiceError(txErr.Error(), "Password Hashing")
		return response, txErr
	}

	userDomain.Password = string(hashPassword)
	
	// 5.4.B Create a new user in the database
	var uuid string
	uuid, txErr = service.UserRepository.Create(ctx, userDomain, tx)
	if txErr != nil {
		return response, txErr
	}

	userDomain.Id = uuid

	// 5.5 Store email verification token to redis 
	emailVerificationDomain := domain.EmailVerification{
		UserId: userDomain.Id,
		Token: emailVerificationToken,
		TTL: 10 * time.Minute,
	}

	service.EmailVerificationRepository.InsertVerification(ctx, &emailVerificationDomain, service.Redis)

	// 6. Create response
	response.UUID = uuid
	response.EmailToken = emailVerificationToken

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
	_, err := service.UserRepository.FetchUserByEmail(ctx, &userDomain, tx)
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

// DONE 12/05/25
func (service *AuthServiceImpl) CheckEmail(request web.UserCheckEmailRequest) (bool, error) {
	// 1. Create necessary variable (e.g context, db)
	var txErr error
	ctx := context.Background()

	tx, txErr := service.DB.Begin()
	if txErr != nil {
		return false, exception.NewServiceError(txErr.Error(), "Transaction Begin")
	}

	// 2. Create defer for db
	defer func() {
		if txErr != nil {
			tx.Rollback()
		} else if recover() != nil {
			tx.Rollback()
		} else { 
			tx.Commit()
		}
	}()

	// 3. Fully validate the request
	err := request.Validate()
	if err != nil {
		txErr = err
		return false, exception.NewServiceError(txErr.Error(), "Request Validation")
	}

	// 4. Bind the data to domain
	userDomain := domain.User{
		Email: request.Email,
	}

	// 5. custom ...
	_, err = service.UserRepository.FetchUserByEmail(ctx, &userDomain, tx)
	if err != nil {
		txErr = err
		if errors.Is(txErr, exception.ErrEmailWasNotFound) {
			return false, nil
		}

		return false, txErr
	}

	// 6. Create response
	return true, nil
}

func (service *AuthServiceImpl) Refresh(user web.UserRefreshRequest) (web.UserRefreshResponse, error) {
	response := web.UserRefreshResponse{}
	tokenDomain := domain.WebToken{
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
	tokenDomain := domain.WebToken{
		RefreshToken: user.RefreshToken,
	}

	ctx := context.Background()
	var txErr error

	// Begin DB Connection
	tx, txErr := service.DB.Begin()
	if txErr != nil {

		return exception.NewServiceError(txErr.Error(), "Begin DB Conn")
	}

	// Create Defer
	defer func() {
		if txErr != nil {
	
			tx.Rollback()	
		} else if recover() != nil {
	
			tx.Rollback()
		} else { 
	
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