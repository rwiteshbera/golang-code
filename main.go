package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Unable to load .env")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "6001"
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	if err := router.Run(":" + PORT); err != nil {
		log.Fatalln("Failed to start server.")
	}
	fmt.Println("Server is listening")
}
