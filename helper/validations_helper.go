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

func GetParsedDateRange(dateFrom, dateTo string) (string, string, error) {
	layout := "2006-01-02"

	// Parse dateFrom
	dateFromParsed, err := time.Parse(layout, dateFrom)
	if err != nil {
		return "", "", fmt.Errorf("error while parsing dateFrom: %v", err)
	}
	// Determine the start of dateFrom and end of dateFrom day
	startOfDateFrom := time.Date(dateFromParsed.Year(), dateFromParsed.Month(), dateFromParsed.Day(), 0, 0, 0, 0, time.UTC)
	endOfDateFrom := startOfDateFrom.Add(24 * time.Hour).Add(-time.Nanosecond)

	var dateToParsed time.Time
	if dateTo == "" || dateTo == dateFrom {
		// If dateTo is not provided or is the same as dateFrom, use the same date as dateTo
		dateToParsed = endOfDateFrom
	} else {
		// Parse dateTo
		dateToParsed, err = time.Parse(layout, dateTo)
		if err != nil {
			return "", "", fmt.Errorf("error while parsing dateTo: %v", err)
		}
		// Ensure dateTo covers the whole day
		startOfDateTo := time.Date(dateToParsed.Year(), dateToParsed.Month(), dateToParsed.Day(), 0, 0, 0, 0, time.UTC)
		dateToParsed = startOfDateTo.Add(24 * time.Hour).Add(-time.Nanosecond)
	}

	// Get the current date
	now := time.Now().UTC()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfToday := startOfToday.Add(24 * time.Hour).Add(-time.Nanosecond)

	// Check if dateFrom is after dateTo
	if dateFromParsed.After(dateToParsed) {
		return "", "", fmt.Errorf("dateFrom should be less than or equal to dateTo")
	}

	// Check if the dates are in the future
	if dateToParsed.After(endOfToday) {
		return "", "", fmt.Errorf("dateFrom and dateTo should be less than or equal to the current date")
	}

	return dateFromParsed.Format(layout), dateToParsed.Format(layout), nil
}
