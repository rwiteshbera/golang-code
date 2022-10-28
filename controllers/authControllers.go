package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
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
		queryForSignup := "INSERT INTO `users` (user_id, username, firstname, lastname, email, hash_password, created_at, last_login) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"

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
			user.LastLogin, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
			// Inserting data
			res, err := statement.ExecContext(ctx, user.UserId, user.UserName, user.FirstName, user.LastName, user.Email, user.Password, user.CreatedAt, user.LastLogin)
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
			c.JSON(http.StatusOK, gin.H{"rows": rows, "username": user.UserName, "fullname": user.FirstName + " " + user.LastName, "email": user.Email, "created_at": user.CreatedAt})
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
		queryForExistingEmailCheck := "SELECT user_id, username, firstname, lastname, email, premium, hash_password FROM users WHERE email = ?"

		var savedUserPassword string

		var savedUserData models.SavedUser

		if err1 := db.QueryRowContext(ctx, queryForExistingEmailCheck, user.Email).Scan(&savedUserData.UserId, &savedUserData.UserName, &savedUserData.FirstName, &savedUserData.LastName, &savedUserData.Email, &savedUserData.IsPremium, &savedUserPassword); err1 == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no user found"})
			return
		}

		// Verify the password
		isPasswordValid, message := helpers.VerifyPassword(user.Password, savedUserPassword)

		if !isPasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": message})
			return
		}

		token, err := helpers.GenerateToken(savedUserData.Email, savedUserData.UserName, savedUserData.FirstName, savedUserData.LastName, savedUserData.UserId, savedUserData.IsPremium)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update the last login data
		sqlQueryToUpdateLastLogin := "UPDATE users SET last_login = ? WHERE email = ?"
		statement, err := db.PrepareContext(ctx, sqlQueryToUpdateLastLogin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		LastLogin, _ := time.Parse(time.RFC1123, time.Now().Format(time.RFC1123)) // When the account is last logged in

		res, err := statement.ExecContext(ctx, LastLogin, savedUserData.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := res.RowsAffected()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token, "rows": rows})
	}
}
