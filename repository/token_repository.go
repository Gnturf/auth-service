package repository

import (
	"auth-service/model/domain"
	"context"
	"database/sql"
)

type TokenRepository interface {
	InsertRefreshToken(ctx context.Context, user domain.User, refreshToken string, tx *sql.Tx) (error)
	FetchRefreshToken(ctx context.Context, token *domain.Tokens, tx *sql.Tx) (domain.Tokens, error)
	RevokeRefreshToken(ctx context.Context, token *domain.Tokens, tx *sql.Tx) (error)
}
