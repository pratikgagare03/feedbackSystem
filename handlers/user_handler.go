package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pratikgagare03/feedback/helper"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var validate = validator.New()

func VerifyPasswordPassword(userPassword, hashedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))

	if err != nil {
		return false, fmt.Sprintf("password is incorrect")
	}

	return true, fmt.Sprintf("")
}
func HashPassword(password string) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hashedPass)
}
func SignUp(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	_, err := repository.GetUserRepository().FindUserByEmail(user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while searching for email."})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists."})
		return
	}
	hashedPass := HashPassword(user.Password)
	user.Password = hashedPass

	// token, refreshToken, err := helper.GenerateAllTokens(user.Email, user.First_name, user.Last_name, user.User_type, user.UserId)
	// if err != nil {
	// 	msg := fmt.Sprintf("failed to generate token")
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
	// 	return
	// }
	// user.Refresh_token = refreshToken

	err = repository.GetUserRepository().InsertUser(&user)
	if err != nil {
		msg := "User item was not created"
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, "user created successfully.")
}
func Login(c *gin.Context) {
	var user models.User
	var foundUser *models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundUser, err := repository.GetUserRepository().FindUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is incorrect"})
		return
	}

	pasIsValid, msg := VerifyPasswordPassword(user.Password, foundUser.Password)
	if !pasIsValid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	if foundUser.Email == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	token, claims, err := helper.GenerateAllTokens(foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.User_type, foundUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token	"})
		return
	}
	c.SetCookie("token", token, int(claims.ExpiresAt), "/", "", true, true)
	// c.SetCookie("refresf_token", refreshToken, int(claims.ExpiresAt), "/", "", true, true)

	helper.UpdateAllTokens(token, foundUser.ID)
	foundUser, err = repository.GetUserRepository().FindUserByID(foundUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foundUser)

}
func GetUsers(c *gin.Context) {
	if err := helper.CheckUserType(c, "ADMIN"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage"))
	if err != nil || recordsPerPage < 1 {
		recordsPerPage = 10
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(c.Query("startIndex"))
	if err != nil || startIndex < 1 {
		startIndex = (page - 1) * recordsPerPage
	}
	users, err := repository.GetUserRepository().GetAllUsersByOffsetLimit(startIndex, recordsPerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	userId := c.Param("user_id")
	if err := helper.MatchUserTypeToUid(c, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user, err := repository.GetUserRepository().FindUserByID(uint(userIdInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
