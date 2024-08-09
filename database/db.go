package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect function is used to connect to the database
func Connect() (*gorm.DB, error){
	host := os.Getenv("HOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	dbname := os.Getenv("DBNAME")

	// Connection string for database creation
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password)

	// Open a connection to the PostgreSQL server
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Logs.Error().Msgf("Error opening connection to PostgreSQL server: %v", err)
		return nil, err
	}
	defer conn.Close()

	// Check if the database exists
	var dbExists bool
	err = conn.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", dbname).Scan(&dbExists)
	if err != nil {
		logger.Logs.Error().Msgf("Error checking if database exists: %v", err)
		return nil, err
	}

	if !dbExists {
		logger.Logs.Info().Msg("Database does not exist. Creating database...")
		_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			logger.Logs.Error().Msgf("Error creating database: %v", err)
			return nil, err
		}
		logger.Logs.Info().Msg("Database created successfully")
	}

	// Connection string for GORM
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open a connection to the specific database with GORM
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
