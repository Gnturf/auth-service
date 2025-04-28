package repository

import (
	"auth-service/model/domain"
	"context"
	"database/sql"
)

type VerificationRepository interface {
	InsertVerification(ctx context.Context, verification *domain.Verification, tx *sql.Tx) (domain.Verification, error)
}