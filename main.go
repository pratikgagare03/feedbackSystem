package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	logger "github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/routes"
)

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
	err := godotenv.Load(".env")
	if err != nil {
		logger.Logs.Error().Msg("Error loading .env file")
	}
	router := gin.Default()
	setupRoutes(router)
	err = router.Run(os.Getenv("APP_PORT"))
	logger.Logs.Error().Msgf("Error starting the server: %v", err)
}
