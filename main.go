package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pratikgagare03/feedback/routes"
)

func setupRoutes(router *gin.Engine) {
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
		log.Fatal("error loading .env file")
	}
	router := gin.Default()
	setupRoutes(router)
	router.Run(os.Getenv("APP_PORT"))
}
