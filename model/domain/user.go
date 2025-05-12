package domain

import "time"

type User struct {
	Id        string
	Email     string
	Password  string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
	EmailVerified bool
	IsActive bool
}