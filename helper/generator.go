package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateRefreshToken(length int) (string, error) {
	return gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", length)
}

func GenerateNumericToken(length int) (string, error) {
	return gonanoid.Generate("0123456789", length)
}

func GenerateJWTAccessToken(id string, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	signedToken, txErr := token.SignedString([]byte(key))

	return signedToken, txErr
}

func XDaysFromNowOnMidnight(x int) time.Time {
	now := time.Now().UTC()
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)


	inXDays := time.Date(todayMidnight.Year(), todayMidnight.Month(), todayMidnight.Day() + x, 23, 59, 59, 59, time.UTC)

	return inXDays
}

func TTLXMin(x int) time.Duration {
	return time.Minute * 2
}

func TTLXDay(x int) time.Time{
	nextXDay := time.Now().AddDate(0, 0, x)
	return nextXDay
}