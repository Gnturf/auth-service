package domain

import "time"

type Verification struct {
	Id        string
	UserId 		string
	Token     string
	ExpiredAt time.Time
	IsUsed    bool
}