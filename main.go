package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pratikgagare03/feedback/handlers"
)

func setupRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/create", handlers.CreateUser)

			feedbackGroup := userGroup.Group("/:userId/feedback")
			{
				feedbackGroup.POST("/create", handlers.CreateFeedback)
				// feedbackGroup.POST("/:{feedbackId}/addQuestion", handlers.AddQuestion)
				feedbackGroup.GET("/:feedbackId", handlers.GetFeedback)
				feedbackGroup.POST("/:feedbackId/respond", handlers.SaveFeedbackResponse)
			}
		}

	}

}
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(".env not found")
	}
	router := gin.Default()
	setupRoutes(router)
	router.Run(os.Getenv("APP_PORT"))
}
