package repository

import (
	"auth-service/model/domain"
	"context"
	"database/sql"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User, tx *sql.Tx) (string, error)
	FetchUserByEmail(ctx context.Context, user *domain.User, tx *sql.Tx) (domain.User, error)
	FetchUserByUsername(ctx context.Context, user *domain.User, tx *sql.Tx) (domain.User, error)
}