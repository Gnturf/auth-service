package repository

import (
	"auth-service/config"
	"auth-service/exception"
	"auth-service/helper"
	"auth-service/model/domain"
	"context"
	"database/sql"
	"log"
)

type TokenRepositoryImpl struct {
	Config *config.AppConfig
	Logger *log.Logger
}

func NewTokenRepository(config *config.AppConfig, logger *log.Logger) TokenRepository {
	return &TokenRepositoryImpl{
		Config: config,
		Logger: logger,
	}
}

func (repository *TokenRepositoryImpl) FetchRefreshToken(ctx context.Context, token *domain.Tokens, tx *sql.Tx) (domain.Tokens, error) {
	// Prepare the SQL query
	query := "SELECT user_id, token FROM refresh_tokens WHERE token = $1 AND expire_at > NOW() AND revoked = false"

	// Fetch to DB
	row := tx.QueryRowContext(ctx, query, token.RefreshToken)

	// Validate
	err := row.Scan(&token.Id, &token.RefreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return *token, exception.NewRepositoryError(exception.RefreshTokenWasNotFound, "Fetch Refresh Token")
		}
		return *token, exception.NewRepositoryError(err.Error(), "Fetch Refresh Token")
	}

	// Return
	return *token, nil
}

func (repository *TokenRepositoryImpl) InsertRefreshToken(ctx context.Context, user domain.User, refreshToken string, tx *sql.Tx) (error) {
	query := "INSERT INTO refresh_tokens (user_id, token, expire_at) VALUES ($1, $2, $3)"

	inXDay := helper.XDaysFromNowOnMidnight(repository.Config.RefreshTokenExpiredDays)

	_, err := tx.ExecContext(ctx, query, user.Id, refreshToken, inXDay)
	if err != nil {
		return exception.NewRepositoryError(err.Error(), "InserRefreshToken")
	}

	return nil
}

func (reqpository *TokenRepositoryImpl) RevokeRefreshToken(ctx context.Context, token *domain.Tokens, tx *sql.Tx) (error) {
	query := "UPDATE refresh_tokens SET revoked = TRUE WHERE token = $1"

	result, err := tx.ExecContext(ctx, query, token.RefreshToken)
	if err != nil {
		return exception.NewRepositoryError(err.Error(), "Query")
	}

	affectedRow, err := result.RowsAffected()
	if err != nil {
		return exception.NewRepositoryError(err.Error(), "Affected Row")
	}

	if affectedRow == 0 {
		return exception.NewRepositoryError(exception.UserWasNotFound, "Affected Row")
	}

	return nil
}