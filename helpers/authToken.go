package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type SignedUserDetails struct {
	UserId    string
	Username  string
	FirstName string
	LastName  string
	Email     string
	IsPremium string
	jwt.RegisteredClaims
}

func GenerateToken(email string, username string, firstname string, lastname string, uid string, ispremium string) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	SECRET_KEY := os.Getenv("SECRET_KEY")

	claims := SignedUserDetails{
		UserId:    uid,
		Username:  username,
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		IsPremium: ispremium,
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

func ValidateToken(authToken string) (claims *SignedUserDetails, message error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	SECRET_KEY := os.Getenv("SECRET_KEY")

	token, err := jwt.ParseWithClaims(
		authToken,
		&SignedUserDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		message = errors.New(err.Error())
	}

	claims, ok := token.Claims.(*SignedUserDetails)

	if !ok {
		message = errors.New("token is invalid")
	}

	return claims, message
}
