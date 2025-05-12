package repository

import (
	"auth-service/config"
	"auth-service/exception"
	"auth-service/model/domain"
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

/**
GENERAL STEP TO WRITE REPOSITORY
1. Write SQL Query
2. Execute Query
3. Custom ...
4. Return the result
*/

type UserRepositoryImpl struct {
	Config *config.AppConfig
}

func NewUserRepositoryImpl(config *config.AppConfig) UserRepository {
	return &UserRepositoryImpl{
		Config: config,
	}
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, user domain.User, tx *sql.Tx) (string, error) {
	var userId string

	query := "INSERT INTO users (email, password_hash, username) VALUES ($1, $2, $3) RETURNING id"

	// Scan the returned id into a variable
	err := tx.QueryRowContext(ctx, query, user.Email, user.Password, user.Username).Scan(&userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return "", exception.NewRepositoryError(exception.UserAlreadyExists, exception.ErrUserAlreadyExist, "Query")
			}
			return "", exception.NewRepositoryError(err.Error(), err,"Query")
		} else {
			return "", exception.NewRepositoryError(err.Error(), err, "Query")
		}
	}

	return userId, nil
}

func (repository *UserRepositoryImpl) FetchUserByEmail(ctx context.Context, user *domain.User, tx *sql.Tx) (domain.User, error) {
	// 1. Write SQL Query
	query := "SELECT id, email, password_hash, username, created_at, updated_at, email_verified, is_active FROM users WHERE email = $1 LIMIT 1"

	// 2. Execute Query
	err := tx.QueryRowContext(ctx, query, user.Email).Scan(&user.Id, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.EmailVerified, &user.IsActive)

	// 3. Custom ...
	// Execute the query with the provided email and transaction context
	if err != nil {
		// If no rows are found (user doesn't exist), return sql.ErrNoRows
		if err == sql.ErrNoRows {
			return *user, exception.NewRepositoryError(exception.EmailWasNotFound, exception.ErrEmailWasNotFound, "Check User Email")
		}
		// Handle other errors that occur during scanning
		return *user, exception.NewRepositoryError(err.Error(), err, "Check User Email")
	}

	// 4. Return the result
	return *user, nil
}

func (repository *UserRepositoryImpl) FetchUserByUsername(ctx context.Context, user *domain.User, tx *sql.Tx) (domain.User, error) {
	// 1. Write SQL Query
	query := "SELECT id, email, password_hash, username, created_at, updated_at, email_verified, is_active FROM users WHERE username = $1 LIMIT 1"

	// 2. Execute Query
	err := tx.QueryRowContext(ctx, query, user.Username).Scan(&user.Id, &user.Email, &user.Password, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.EmailVerified, &user.IsActive)

	// 3. Custom ...
	// Execute the query with the provided email and transaction context
	if err != nil {
		// If no rows are found (user doesn't exist), return sql.ErrNoRows
		if err == sql.ErrNoRows {
			return *user, exception.NewRepositoryError(exception.UsernameWasNotFound, exception.ErrUsernameWasNotFound, "Check User Username")
		}
		// Handle other errors that occur during scanning
		return *user, exception.NewRepositoryError(err.Error(), err, "Check User Username")
	}

	// 4. Return the result
	return *user, nil
}