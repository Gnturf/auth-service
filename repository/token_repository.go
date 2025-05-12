package repository

import (
	"auth-service/model/domain"
	"context"
	"database/sql"
)

type TokenRepository interface {
	InsertRefreshToken(ctx context.Context, user domain.User, refreshToken string, tx *sql.Tx) (error)
	FetchRefreshToken(ctx context.Context, token *domain.WebToken, tx *sql.Tx) (domain.WebToken, error)
	RevokeRefreshToken(ctx context.Context, token *domain.WebToken, tx *sql.Tx) (error)
}
