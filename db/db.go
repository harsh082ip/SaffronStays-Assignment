package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitDB() (*sql.DB, error) {

	godotenv.Load(".env")

	serviceURI := os.Getenv("SERVICE_URI")
	if serviceURI == "" {
		log.Fatal("SERVICE_URI is not set. Please provide the database connection string.")
	}

	db, err := sql.Open("postgres", serviceURI)
	if err != nil {
		return &sql.DB{}, err
	}

	return db, nil
}
