package domain

import "time"

type FullUser struct {
	UserProfile
	PasswordHash string
}

type UserProfile struct {
	ID         string
	FirstName  string
	SecondName string
	Birthdate  time.Time
	Biography  string
	City       string
}
