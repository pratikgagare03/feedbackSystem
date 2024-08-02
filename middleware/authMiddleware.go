package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/helper"
	"github.com/pratikgagare03/feedback/logger"
)

func Authenticate(c *gin.Context) {
	clientToken, err := c.Cookie("token")
	if err != nil {
		logger.Logs.Error().Msgf("error while getting cookie: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No cookie found, please login"})
		c.Abort()
		return
	}

	claims, msg := helper.ValidateToken(clientToken)
	if msg != "" {
		logger.Logs.Error().Msgf("error while validating token: %v", msg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		c.Abort()
		return
	}
	c.Set("email", claims.Email)
	c.Set("first_name", claims.First_name)
	c.Set("last_name", claims.Last_name)
	c.Set("uid", claims.Uid)
	c.Set("user_type", claims.User_type)
	c.Next()
}
