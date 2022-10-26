package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		return "", errors.New("unable to hash the password")
	}
	return string(hashedPassword), nil
}

// Verify the password with saved user password
func VerifyPassword(userPassword string, hashedPassword string) (bool, string) {
	check := true
	message := ""

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword)); err != nil {
		message = "password is incorrect"
		check = false
	}

	return check, message
}
