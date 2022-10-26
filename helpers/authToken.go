package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type SignedUserDetails struct {
	Email     string
	FirstName string
	LastName  string
	UserId    string
	jwt.RegisteredClaims
}

func GenerateToken(email string, firstname string, lastname string, uid string) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	SECRET_KEY := os.Getenv("SECRET_KEY")

	claims := SignedUserDetails{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		UserId:    uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 168)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		fmt.Println(err)
	}
	return token, err
}
