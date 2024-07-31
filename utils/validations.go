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
	if err != nil || err == gorm.ErrRecordNotFound {
		return false, err
	}

	return true, nil
}

func IsValidFeedbackId(feedbackId string) (bool, error) {
	_, err := repository.GetFeedbackRepository().FindFeedbackByID(feedbackId)
	if err != nil || err == gorm.ErrRecordNotFound {
		return false, err
	}
	return true, nil
}

func ResponseExistForUser(feedbackID string, userID string) bool {
	responses, err := repository.GetResponseRepository().FindResponseByUserIdFeedbackId(userID, feedbackID)
	if len(responses) != 0 && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func QuestionExistInFeedback(questionID string, feedbackID string) bool {
	res, err := repository.GetQuestionRepository().FindQuestionByQuestionIdFeedbackId(questionID, feedbackID)
	if len(res) != 0 && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func ResponseExistForFeedback(feedbackID string) bool {
	responses, err := repository.GetResponseRepository().FindResponseByFeedbackId(feedbackID)
	if len(responses) != 0 && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}
