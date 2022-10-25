package models

import "time"

type User struct {
	UserId    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"hash_password"`
	CreatedAt time.Time `json:"created_at"`
}
