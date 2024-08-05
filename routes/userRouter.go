package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/handlers"
	"github.com/pratikgagare03/feedback/middleware"
)

func UserRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.Use(middleware.Authenticate)
	incomingRoutes.GET("/users", handlers.GetUsers)
	incomingRoutes.GET("/users/:user_id", handlers.GetUser)
	incomingRoutes.GET("/user/responses", handlers.GetAllResponsesForUser)
}
