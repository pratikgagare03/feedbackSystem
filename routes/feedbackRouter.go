package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/handlers"
	"github.com/pratikgagare03/feedback/middleware"
)

func FeedbackRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.Use(middleware.Authenticate)
	incomingRoutes.POST("/create", handlers.CreateFeedback)
	incomingRoutes.GET("/:feedbackId", handlers.GetFeedback)
	incomingRoutes.POST("/:feedbackId/respond", handlers.SaveFeedbackResponse)
	incomingRoutes.GET("/:feedbackId/responses", handlers.GetAllResponsesForFeedback)
	incomingRoutes.PATCH("/:feedbackId/publish", handlers.TogglePublishStatus(true))
	incomingRoutes.PATCH("/:feedbackId/unpublish", handlers.TogglePublishStatus(false))
	incomingRoutes.GET("/:feedbackId/stats", handlers.GetFeedbackStats)
}
