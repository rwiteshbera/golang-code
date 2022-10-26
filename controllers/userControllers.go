package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type savedUser struct {
	UserName string
	FullName string
	Email    string
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("name")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		queryForUserDataFetch := "SELECT username, fullname, email FROM users WHERE username = ?"

		var userdata savedUser

		if err := db.QueryRowContext(ctx, queryForUserDataFetch, username).Scan(&userdata.UserName, &userdata.FullName, &userdata.Email); err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no user found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": userdata})
	}
}
