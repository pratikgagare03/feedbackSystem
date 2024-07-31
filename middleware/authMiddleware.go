package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/helper"
)

func Authenticate(c *gin.Context) {
	clientToken, err := c.Cookie("token")
	if err!=nil ||clientToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No cookie found, please login"})
		c.Abort()
		return
	}

	claims, msg := helper.ValidateToken(clientToken)
	if msg != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
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
