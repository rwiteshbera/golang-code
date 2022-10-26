package models

import "time"

type User struct {
	UserId    string    `json:"user_id"`
	UserName  string    `json:"username"`
	FullName  string    `json:"fullname"`
	Email     string    `json:"email"`
	Password  string    `json:"hash_password"`
	CreatedAt time.Time `json:"created_at"`
}
