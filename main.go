package main

import (
	"os"

	"github.com/gin-gonic/gin"
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
	logger.Logs.Info().Msg("Starting the server")
	router := gin.Default()
	setupRoutes(router)
	logger.Logs.Info().Msgf("Starting the server on port %s", os.Getenv("APP_PORT"))
	err := router.Run(os.Getenv("APP_PORT"))
	logger.Logs.Fatal().Msgf("Error starting the server: %v", err)
}
