package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/utils"
)

func SaveFeedbackResponse(c *gin.Context) {
	userId := c.Param("userId")
	if ok, err := utils.IsValidUser(userId); !ok {
		log.Printf("ERROR:invalid userId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var responseInput models.FeedbackResponse
	if err := c.ShouldBindJSON(&responseInput); err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := validate.Struct(responseInput)
	if err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, responseInput)
}
