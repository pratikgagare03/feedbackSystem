package models

import (
	"gorm.io/gorm"
)

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
type FeedbackResponseInput struct {
	gorm.Model
	QuestionAnswer []QuestionAnswer
}
type QuestionAnswer struct {
	gorm.Model
	QuestionID uint   `validate:"required"`
	Answer     string `validate:"required"`
}
