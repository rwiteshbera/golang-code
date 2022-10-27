package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rwiteshbera/orbit/config"
	"github.com/rwiteshbera/orbit/helpers"
	"github.com/rwiteshbera/orbit/models"
)

var db *sql.DB

func init() {
	db = config.GetDatabase()
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		// Storing request body in user of type User struct
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Query for inserting data into users table
		queryForSignup := "INSERT INTO `users` (user_id, username, fullname, email, hash_password, created_at) VALUES(?, ?, ?, ?, ?, ?)"

		// Query to fetch row by matching the email // Check for existing user
		queryForExistingEmailOrUsernameCheck := "SELECT * FROM users WHERE (email = ? OR username = ?)"

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var userid string

		// Query for user
		err1 := db.QueryRowContext(ctx, queryForExistingEmailOrUsernameCheck, user.Email, user.UserName).Scan(&userid)

		// If user doesn't exist then do signup
		if err1 == sql.ErrNoRows {
			statement, err := db.PrepareContext(ctx, queryForSignup)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
			defer statement.Close()

			hashedPassword, err := helpers.HashPassword(user.Password) // Hash the password using bcrypt
			if err != nil {
				fmt.Println("error: " + err.Error())
				return
			}

			user.UserId = uuid.New().String() // Create a random id for user
			user.Password = hashedPassword
			user.CreatedAt, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123)) // When the account is created

			// Inserting data
			res, err := statement.ExecContext(ctx, user.UserId, user.UserName, user.FullName, user.Email, user.Password, user.CreatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}

			// rows affected == signup successful;
			rows, err := res.RowsAffected()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}

			// return a response
			c.JSON(http.StatusOK, gin.H{"rows": rows, "username": user.UserName, "fullname": user.FullName, "email": user.Email, "created_at": user.CreatedAt})
			return
		} else {
			// If the email is already used by user
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Email or Username already exists!"})
			return
		}

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User // Requested user data

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Query to fetch row by matching the email // Check for existing user
		queryForExistingEmailCheck := "SELECT user_id, username, fullname, email, premium, hash_password FROM users WHERE email = ?"

		var savedUserPassword string

		var userdata models.SavedUser

		if err1 := db.QueryRowContext(ctx, queryForExistingEmailCheck, user.Email).Scan(&userdata.UserId, &userdata.UserName, &userdata.FullName, &userdata.Email, &userdata.IsPremium, &savedUserPassword); err1 == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no user found"})
			return
		}

		// Verify the password
		isPasswordValid, message := helpers.VerifyPassword(user.Password, savedUserPassword)

		if !isPasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": message})
			return
		}

		// Split the firstname and lastname
		splitName := strings.Split(userdata.FullName, " ")
		var savedUserFirstName string = splitName[0]
		var savedUserLastName string = splitName[1]

		token, err := helpers.GenerateToken(userdata.Email, userdata.UserName, savedUserFirstName, savedUserLastName, userdata.UserId, userdata.IsPremium)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
