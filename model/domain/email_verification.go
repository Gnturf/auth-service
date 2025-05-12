package domain

import "time"

type EmailVerification struct {
	UserId 		string
	Token     string
	TTL				time.Duration
}