package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/handlers"
)

func AuthRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("users/signup", handlers.SignUp)
	incomingRoutes.POST("users/login", handlers.Login)
}
