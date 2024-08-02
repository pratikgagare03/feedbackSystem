package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pratikgagare03/feedback/helper"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var validate = validator.New()

// VerifyPasswordPassword verifies the password
func VerifyPasswordPassword(userPassword, hashedPassword string) (bool, string) {
	logger.Logs.Info().Msg("verifying password")
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))

	if err != nil {
		logger.Logs.Error().Msgf("error while verifying password: %v", err)
		return false, "password is incorrect"
	}

	return true, ""
}
func HashPassword(password string) string {
	logger.Logs.Info().Msg("hashing password")
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		logger.Logs.Error().Msgf("error while hashing password: %v", err)
	}
	return string(hashedPass)
}
func SignUp(c *gin.Context) {
	logger.Logs.Info().Msg("signing up")
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		logger.Logs.Error().Msgf("error while binding json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		logger.Logs.Error().Msgf("error while validating user: %v", validationErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	_, err := repository.GetUserRepository().FindUserByEmail(user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Logs.Error().Msgf("error while searching for email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while searching for email."})
		return
	} else if err != gorm.ErrRecordNotFound {
		logger.Logs.Error().Msg("email already exists")
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
		logger.Logs.Error().Msgf("error while inserting user: %v", err)
		msg := "User item was not created"
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	logger.Logs.Info().Msg("user created successfully")
	c.JSON(http.StatusOK, "user created successfully.")
}
func Login(c *gin.Context) {
	logger.Logs.Info().Msg("logging in")
	var user models.User
	var foundUser *models.User

	if err := c.BindJSON(&user); err != nil {
		logger.Logs.Error().Msgf("error while binding json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundUser, err := repository.GetUserRepository().FindUserByEmail(user.Email)
	if err != nil {
		logger.Logs.Error().Msgf("error while searching for email: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is incorrect"})
		return
	}

	pasIsValid, msg := VerifyPasswordPassword(user.Password, foundUser.Password)
	if !pasIsValid {
		logger.Logs.Error().Msgf("error while verifying password: %v", msg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	if foundUser.Email == "" {
		logger.Logs.Error().Msg("user not found")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}
	token, claims, err := helper.GenerateAccessToken(foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.User_type, foundUser.ID)
	if err != nil {
		logger.Logs.Error().Msgf("error while generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generating token	"})
		return
	}
	c.SetCookie("token", token, int(claims.ExpiresAt), "/", "", true, true)

	foundUser, err = repository.GetUserRepository().FindUserByID(foundUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logs.Info().Msg("logged in successfully")
	c.JSON(http.StatusOK, foundUser)

}
func GetUsers(c *gin.Context) {
	logger.Logs.Info().Msg("getting users")
	if err := helper.CheckUserType(c, "ADMIN"); err != nil {
		logger.Logs.Error().Msgf("error while checking user type: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage"))
	if err != nil || recordsPerPage < 1 {
		logger.Logs.Error().Msg("error while getting records per page")
		logger.Logs.Info().Msg("setting records per page to 10")
		recordsPerPage = 10
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		logger.Logs.Error().Msg("error while getting page")
		logger.Logs.Info().Msg("setting page to 1")
		page = 1
	}

	startIndex, err := strconv.Atoi(c.Query("startIndex"))
	if err != nil || startIndex < 1 {
		logger.Logs.Error().Msg("error while getting start index")
		startIndex = (page - 1) * recordsPerPage
	}
	users, err := repository.GetUserRepository().GetAllUsersByOffsetLimit(startIndex, recordsPerPage)
	if err != nil {
		logger.Logs.Error().Msgf("error while listing user items: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		return
	}
	logger.Logs.Info().Msg("users listed successfully")
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	userId := c.Param("user_id")
	if err := helper.MatchUserTypeToUid(c, userId); err != nil {
		logger.Logs.Error().Msgf("error while matching user type to uid: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		logger.Logs.Error().Msgf("error while converting user id to int: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user, err := repository.GetUserRepository().FindUserByID(uint(userIdInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Logs.Error().Msg("user not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Logs.Info().Msg("user found successfully")
	c.JSON(http.StatusOK, user)
}
