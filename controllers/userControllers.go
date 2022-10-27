package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/orbit/models"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		loggedInUsername, _ := c.Get("username") // Fetching authorized username

		username := c.DefaultQuery("name", loggedInUsername.(string)) // If no username is provided, then it will return authorized user's own details

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		queryForUserDataFetch := "SELECT username, fullname, email, premium FROM users WHERE username = ?"

		var userdata models.SavedUser

		if err := db.QueryRowContext(ctx, queryForUserDataFetch, username).Scan(&userdata.UserName, &userdata.FullName, &userdata.Email, &userdata.IsPremium); err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no user found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"username": userdata.UserName, "fullname": userdata.FullName, "email": userdata.Email, "premium": userdata.IsPremium})
	}
}
