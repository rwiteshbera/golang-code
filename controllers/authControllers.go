package controllers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rwiteshbera/orbit/config"
	"github.com/rwiteshbera/orbit/models"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func init() {
	db = config.GetDatabase()
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		return "", errors.New("unable to hash the password")
	}
	return string(hashedPassword), nil
}

func Signup(c *gin.Context) {
	var user models.User

	err := c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO `users` (user_id, name, email, hash_password, created_at) VALUES(?, ?, ?, ?, ?)"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	statement, err := db.PrepareContext(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	defer statement.Close()

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		fmt.Println("error: " + err.Error())
	}

	user.UserId = uuid.New().String()
	user.Password = hashedPassword
	user.CreatedAt, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))

	res, err := statement.ExecContext(ctx, user.UserId, user.Name, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": rows})

}

func Login() {

}
