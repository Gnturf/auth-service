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

type VerificationRepositoryImpl struct {
	Config *config.AppConfig
	Logger *log.Logger
}

func NewVerificationRepository() VerificationRepository {
	return &VerificationRepositoryImpl{}
}

func (r *VerificationRepositoryImpl) InsertVerification(ctx context.Context, verification *domain.Verification, tx *sql.Tx) (domain.Verification, error) {
	query := "INSERT INTO email_verification(user_id, token, expire_at) VALUES(?, ?, ?)"

	helper.XDaysFromNowOnMidnight()

	res, err := tx.ExecContext(ctx, query, verification.UserId, verification.Token)
	if err != nil {
		return *verification, exception.NewRepositoryError(err.Error(), "Query")
	}

}