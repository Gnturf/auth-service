package repository

import (
	"auth-service/model/domain"
	"context"

	"github.com/redis/go-redis/v9"
)

type EmailVerificationRepository interface {
	InsertVerification(ctx context.Context, verification *domain.EmailVerification, client *redis.Client) (domain.EmailVerification, error)
}