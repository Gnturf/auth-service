package repository

import (
	"auth-service/config"
	"auth-service/exception"
	"auth-service/model/domain"
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepositoryImpl struct {
	Config *config.AppConfig
	Logger *log.Logger
}

func NewAuthRepository(config *config.AppConfig, logger *log.Logger) AuthRepository {
	return &AuthRepositoryImpl{
		Config: config,
		Logger: logger,
	}
}

func (repository *AuthRepositoryImpl) CreateUser(ctx context.Context, user domain.User, tx *sql.Tx) (string, error) {
	var userId string

	query := "INSERT INTO users (email, password_hash, username) VALUES ($1, $2, $3) RETURNING id"

	// Scan the returned id into a variable
	err := tx.QueryRowContext(ctx, query, user.Email, user.Password, user.Username).Scan(&userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				// Unique constraint violation
				repository.Logger.Printf("[%s] User already exists: %s", user.Id, pgErr.Message)
				return "", exception.NewRepositoryError(exception.UserAlreadyExists, "Query")
			}
			repository.Logger.Printf("[%s] PGX Query Error: %s", user.Id, err.Error())
			return "", exception.NewRepositoryError(err.Error(), "Query")
		} else {
			repository.Logger.Printf("[%s] Query Error: %s", user.Id, err.Error())
			return "", exception.NewRepositoryError(err.Error(), "Query")
		}
	}

	return userId, nil
}

func (repository *AuthRepositoryImpl) FetchUser(ctx context.Context, user *domain.User, tx *sql.Tx) (domain.User, error) {

	// Prepare the SQL query to select the hashed password for the given email
	query := "SELECT id, email, password_hash, username FROM users WHERE email = $1 LIMIT 1"

	// Execute the query with the provided email and transaction context
	err := tx.QueryRowContext(ctx, query, user.Email).Scan(&user.Id, &user.Email, &user.Password, &user.Username)
	if err != nil {
		// If no rows are found (user doesn't exist), return sql.ErrNoRows
		if err == sql.ErrNoRows {
			return *user, exception.NewRepositoryError(exception.UserWasNotFound, "CheckUser")
		}
		// Handle other errors that occur during scanning
		return *user, exception.NewRepositoryError(err.Error(), "CheckUser")
	}

	// Return the populated user struct
	return *user, nil
}

