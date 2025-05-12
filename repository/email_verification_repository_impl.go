package repository

import (
	"auth-service/config"
	"auth-service/exception"
	"auth-service/model/domain"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type VerificationRepositoryImpl struct {
	Config *config.AppConfig
}

func NewVerificationRepositoryImpl() EmailVerificationRepository {
	return &VerificationRepositoryImpl{}
}

func (r *VerificationRepositoryImpl) InsertVerification(ctx context.Context, verification *domain.EmailVerification, client *redis.Client) (domain.EmailVerification, error) {
	// Creating the key; format <service>:<domain>:<user_id>
	key := fmt.Sprintf("auth:email_verification:%s", verification.UserId)

	status := client.Set(ctx, key, verification.Token, verification.TTL)

	// Check if the operation was successful
	if status.Err() != nil {
			return *verification, exception.NewRepositoryError(status.Err().Error(), status.Err(), "Set Redis Key")
	} else {
			return *verification, nil
	}
}