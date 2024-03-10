package domain

import "time"

type FullUser struct {
	ID string
	UserProfile
	PasswordHash string
}

type UserProfile struct {
	FirstName  string
	SecondName string
	Birthdate  time.Time
	Biography  string
	City       string
}
