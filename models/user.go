package models

import "time"

type User struct {
	UserId    string    `json:"user_id"`
	UserName  string    `json:"username"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"hash_password"`
	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login"`
}

type SavedUser struct {
	UserId    string
	UserName  string
	FirstName string
	LastName  string
	Email     string
	IsPremium string
	LastLogin string
}
