package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rwiteshbera/orbit/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("token")

		if authToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "you must be logged in to the server (unauthorized).",
			})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(authToken)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
		}

		c.Set("userid", claims.UserId)
		c.Next()
	}
}
