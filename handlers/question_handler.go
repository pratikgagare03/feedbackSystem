package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/models"
)
func CreateQuestion(c *gin.Context) {
	var que models.Question
	c.JSON(http.StatusCreated, que)
}

func GetquestionHandler(c *gin.Context) {
	var que models.Question
	c.JSON(http.StatusOK, que)
}

// UpdatequestionHandler handles the PUT request to update an existing question
func UpdatequestionHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "question updated successfully"})
}

// DeletequestionHandler handles the DELETE request to delete a specific question
func DeletequestionHandler(c *gin.Context) {

	c.Status(http.StatusNoContent)
}
