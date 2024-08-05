package helper

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/repository"
	"gorm.io/gorm"
)

// CheckUserType checks if the user type is the same as the role (admin, user, etc)
func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	if userType != role {
		err := errors.New("unauthorized to access this resource")
		return err
	}
	return nil
}

// MatchUserTypeToUid checks if the user type matches the user id
func MatchUserTypeToUid(c *gin.Context, userId string) error {
	userType := c.GetString("user_type")
	uid := c.GetUint("uid")
	// if the user type is not admin, then the user id should match the user id in the token i.e the user should be able to access only his data
	if userType == "USER" && strconv.Itoa(int(uid)) != userId {
		logger.Logs.Error().Msg("Unauthorized to access this resource")
		err := errors.New("unauthorized to access this resource")
		return err
	}

	return nil
}

func MatchFeedbackOwner(c *gin.Context, feedbackID string) (bool, error) {
	feedback, err := repository.GetFeedbackRepository().FindFeedbackByID(feedbackID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	} else if feedback.UserID == c.GetUint("uid") {
		return true, nil
	}
	return false, nil
}
