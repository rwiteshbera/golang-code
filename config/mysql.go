package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func GetDatabase() *sql.DB {
	DB = ConnectDB()
	return DB
}

func ConnectDB() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error while loading .env in /config/mysql.go")
	}

	dsn := os.Getenv("DB_DSN") // Database Data Source Name

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Unable to connect with mysql database")
	}
	fmt.Println("Connected to MySQL Database")

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()
	return db
}
