package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pratikgagare03/feedback/handlers"
	"github.com/pratikgagare03/feedback/logger"
)

func setupRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/create", handlers.CreateUser)
		}
	}

}
func main() {
	logger.Logs.Info().Msg("Started Main")
	err := godotenv.Load(".env")
	if err != nil {
		logger.Logs.Error().Err(err)
	}
	defer logger.File.Close()
	router := gin.Default()
	setupRoutes(router)
	router.Run(os.Getenv("APP_PORT"))
	logger.Logs.Info().Msg("Main Function over")
}
