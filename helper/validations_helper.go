package helper

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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
	if feedbackId == "" {
		return false, errors.New("feedbackId is empty")
	}
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

func IsFeedbackPublished(feedbackID string) (bool, error) {
	feedback, err := repository.GetFeedbackRepository().FindFeedbackByID(feedbackID)
	if err != nil {
		return false, err
	}
	return feedback.Published, nil
}

func GetParsedDateRange(dateFrom, dateTo string) (time.Time, time.Time, error) {
	dateFromParsed, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error while parsing dateFrom: %v", err)
	}

	dateToParsed, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error while parsing dateTo: %v", err)
	}

	if dateFromParsed.Format("2006-01-02") > dateToParsed.Format("2006-01-02") {
		return time.Time{}, time.Time{}, fmt.Errorf("dateFrom should be less than dateTo")
	}

	if dateFromParsed.Format("2006-01-02") > time.Now().Format("2006-01-02") || dateToParsed.Format("2006-01-02") > time.Now().Format("2006-01-02") {
		return time.Time{}, time.Time{}, fmt.Errorf("dateFrom and dateTo should be less than or equal to current date")
	}

	return dateFromParsed, dateToParsed, nil
}
