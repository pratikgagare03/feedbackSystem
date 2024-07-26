package helper

import (
	"errors"

	"github.com/pratikgagare03/feedback/models"
)

func GetQuestionType(s string) (models.QuestionType, error) {
	switch s {
	case "mcq":
		return models.MCQ, nil
	case "textinput":
		return models.TextInput, nil
	case "ratings":
		return models.Ratings, nil
	default:
		return "", errors.New("invalid question type, required enum(\"mcq\",\"textinput\",\"ratings\")")

	}
}

// func GetQueInputToQueModel(queInput []models.QuestionInput) ([]models.Question, error){

// }
