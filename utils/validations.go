package utils

import (
	"strconv"

	"github.com/pratikgagare03/feedback/repository"
	"gorm.io/gorm"
)

func IsValidUser(userId string) (bool, error) {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return false, err
	}
	_, err = repository.GetUserRepository().FindUserByID(uint(userIdInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func IsValidFeedbackId(feedbackId string) (bool, error) {
	_, err := repository.GetFeedbackRepository().FindFeedbackByID(feedbackId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ResponseExistForUser(feedbackID string, userID uint) (bool, error) {
	_, err := repository.GetResponseRepository().FindResponseByUserIdFeedbackId(userID, feedbackID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func QuestionExistInFeedback(questionID uint, feedbackID string) (bool, error) {
	_, err := repository.GetQuestionRepository().FindQuestionByQuestionIdFeedbackId(questionID, feedbackID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ResponseExistForFeedback(feedbackID string) (bool, error) {
	_, err := repository.GetResponseRepository().FindResponseByFeedbackId(feedbackID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
