package config

import (
	"os"
	"strconv"
)

type AppConfig struct {
	Port                    string
	SecretKey               string
	RefreshTokenExpiredDays int
	RefreshTokenLength      int
	BcryptHashCost 					int
}

var AppConfigInstance AppConfig

func LoadEnv() {
	AppConfigInstance.Port = os.Getenv("PORT")
	AppConfigInstance.SecretKey = os.Getenv("SECRET_KEY")
	AppConfigInstance.RefreshTokenExpiredDays = LoadEnvInt("REFRESH_TOKEN_EXPIRED_DAYS")
	AppConfigInstance.RefreshTokenLength = LoadEnvInt("REFRESH_TOKEN_LENGTH")
	AppConfigInstance.BcryptHashCost = LoadEnvInt("BCRYPT_HASH_COST")
}

func LoadEnvInt(key string) int {
	tokenExpiredAt := os.Getenv(key)

	expiredAt, err := strconv.Atoi(tokenExpiredAt)
	if err != nil {
		panic(err)
	}

	return expiredAt
}
