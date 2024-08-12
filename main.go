package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	logger "github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/repository"
	"github.com/pratikgagare03/feedback/routes"
)

func init() {
	logger.StartLogger()
	err := godotenv.Load(".env")
	if err != nil {
		logger.Logs.Fatal().Msgf("Error loading .env file: %v", err)
	}
	err = repository.Connect()
	if err != nil {
		logger.Logs.Fatal().Msgf("Error connecting to the database: %v", err)
	}
}
func setupRoutes(router *gin.Engine) {
	logger.Logs.Info().Msg("Setting up api routes")
	apiGroup := router.Group("/api")
	{
		routes.AuthRoutes(apiGroup)
		routes.UserRoutes(apiGroup)

		feedbackGroup := apiGroup.Group("/feedback")
		routes.FeedbackRoutes(feedbackGroup)
	}

}
func main() {
	logger.Logs.Info().Msg("Starting the server")
	router := gin.Default()
	setupRoutes(router)
	logger.Logs.Info().Msgf("Starting the server on port %s", os.Getenv("APP_PORT"))
	err := router.Run(os.Getenv("APP_PORT"))
	logger.Logs.Fatal().Msgf("Error starting the server: %v", err)
}
