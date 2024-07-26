package handlers

import (
	"context"
	"log"
	"net/http"
	"net/mail"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
)

var validate = validator.New()

func CreateUser(c *gin.Context) {
	var User models.User
	if err := c.ShouldBindJSON(&User); err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := validate.Struct(User)
	if err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if User.Email == "" || User.Password == "" {
		log.Printf("ERROR %+v", "empty email or password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty email or password"})
		return
	}
	_, err = mail.ParseAddress(User.Email)
	if err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if User.Username == "" {
		User.Username = User.Email
		log.Print("No username provided, Used email as username")
	}

	err = repository.GetUserRepository().InsertUser(context.TODO(), &User)

	if err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, User)
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
