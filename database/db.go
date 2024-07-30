package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setup() (*gorm.DB, error) {
	godotenv.Load(".env")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("HOST"), os.Getenv("DBPORT"), os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBNAME"))
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	db.AutoMigrate(models.User{}, models.Feedback{}, models.Question{}, models.FeedbackResponse{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
