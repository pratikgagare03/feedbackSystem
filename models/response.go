package models

import (
	"gorm.io/gorm"
)// FeedbackResponse table in database.
type FeedbackResponse struct {
	gorm.Model
	UserID          uint
	User            User `gorm:"foreignKey:UserID" json:"-"`
	FeedbackID      uint
	Feedback        Feedback `gorm:"foreignKey:FeedbackID" json:"-"`
	QuestionContent string
	QuestionType    QuestionType
	Answer          string
}

// FeedbackResponseInput is a struct that represents the input for feedback response.
type FeedbackResponseInput struct {
	gorm.Model
	QuestionAnswer []QuestionAnswer
}

// QuestionAnswer is a struct that represents the answer for a question.
type QuestionAnswer struct {
	gorm.Model
	QuestionID uint   `validate:"required"`
	Answer     string `validate:"required"`
}
