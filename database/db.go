package database

import (
	"fmt"
	"os"

	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect function is used to connect to the database
func Connect() (*gorm.DB, error) {
	logger.Logs.Info().Msg("Creating connection string")
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("HOST"), os.Getenv("DBPORT"), os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBNAME"))

	logger.Logs.Info().Msg("Connecting to the database")
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		logger.Logs.Error().Msgf("Error connecting to the database: %v", err)
		return nil, err
	}

	logger.Logs.Info().Msg("Running database migrations")
	err = db.AutoMigrate(models.User{}, models.Feedback{}, models.Question{}, models.FeedbackResponse{}, models.Options{}, models.RatingsRange{})
	if err != nil {
		logger.Logs.Error().Msgf("Error running migrations: %v", err)
		return nil, err
	}

	logger.Logs.Info().Msg("Database connection established successfully")
	return db, nil
}
