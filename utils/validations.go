package utils

import (
	"context"

	"github.com/pratikgagare03/feedback/repository"
	"gorm.io/gorm"
)

func IsValidUser(userId string) (bool, error) {
	_, err := repository.GetUserRepository().FindUserByID(context.TODO(), userId)
	if err != nil || err == gorm.ErrRecordNotFound {
		return false, err
	}
	
	return true, nil
}

func IsValidFeedbackId(feedbackId string) (bool, error) {
	_, err := repository.GetFeedbackRepository().FindFeedbackByID(context.TODO(), feedbackId)
	if err != nil || err == gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

