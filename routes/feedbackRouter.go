package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/handlers"
	"github.com/pratikgagare03/feedback/middleware"
)

func FeedbackRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.Use(middleware.Authenticate)
	incomingRoutes.POST("/create", handlers.CreateFeedback)
	// feedbackGroup.POST("/:{feedbackId}/addQuestion", handlers.AddQuestion)
	incomingRoutes.GET("/:feedbackId", handlers.GetFeedback)
	incomingRoutes.POST("/respond/:feedbackId", handlers.SaveFeedbackResponse)
	incomingRoutes.GET("/respond/:userID", handlers.GetAllResponsesForUser)
}
