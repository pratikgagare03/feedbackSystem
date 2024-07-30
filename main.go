package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pratikgagare03/feedback/handlers"
	"github.com/pratikgagare03/feedback/routes"
)

func setupRoutes(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/signup", handlers.SignUp)
		apiGroup.POST("/login", handlers.Login)

		feedbackGroup := apiGroup.Group("/feedback")
		{
			feedbackGroup.POST("/create", handlers.CreateFeedback)
			// feedbackGroup.POST("/:{feedbackId}/addQuestion", handlers.AddQuestion)
			feedbackGroup.GET("/:feedbackId", handlers.GetFeedback)
			feedbackGroup.POST("/:feedbackId/respond", handlers.SaveFeedbackResponse)
		}

	}

}
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(".env not found")
	}
	router := gin.Default()
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	setupRoutes(router)

	router.Run(os.Getenv("APP_PORT"))
}
