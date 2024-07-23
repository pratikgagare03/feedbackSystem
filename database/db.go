package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setup() (*gorm.DB, error) {
	godotenv.Load(".env")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", os.Getenv("Host"), os.Getenv("Port"), os.Getenv("User"), os.Getenv("Name"))

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
