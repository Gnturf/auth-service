package repository

import (
	"auth-service/model/domain"
	"context"
	"database/sql"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user domain.User, tx *sql.Tx) (string, error)
	FetchUser(ctx context.Context, user *domain.User, tx *sql.Tx) (domain.User, error)
}