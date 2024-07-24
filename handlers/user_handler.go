package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
)

func CreateUser(c *gin.Context) {
	var UserInput models.UserInput
	if err := c.ShouldBindJSON(&UserInput); err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	repository.GetUserRepository().InsertUser(context.TODO(), &UserInput)
	c.JSON(http.StatusCreated, gin.H{"message":"user created successfully"})
}

func GetuserHandler(c *gin.Context) {
	var que models.User
	c.JSON(http.StatusOK, que)
}

// UpdateuserHandler handles the PUT request to update an existing user
func UpdateuserHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

// DeleteuserHandler handles the DELETE request to delete a specific user
func DeleteuserHandler(c *gin.Context) {

	c.Status(http.StatusNoContent)
}
