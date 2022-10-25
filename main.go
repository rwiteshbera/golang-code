package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	config "github.com/rwiteshbera/orbit/config"
)

var db *sql.DB

func init() {
	db = config.GetDatabase()
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Unable to load .env in main.go")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "6001"
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// routes.AuthRoutes(router) // routes/authRoutes.go
	// routes.UserRoutes(router) // routes/userRoutes.go

	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	if err := router.Run(":" + PORT); err != nil {
		log.Fatalln("Failed to start server.")
	}
	fmt.Println("Server is listening")
}
