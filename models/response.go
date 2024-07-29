package models

import (
	"gorm.io/gorm"
)

type FeedbackResponse struct {
	gorm.Model
	UserID     uint
	User       User `gorm:"foreignKey:UserID" json:"-"`
	FeedbackID uint
	Feedback   Feedback `gorm:"foreignKey:FeedbackID" json:"-"`
	QuestionID string
	Answer     string
}
type FeedbackResponseInput struct {
	gorm.Model
	UserID         uint
	User           User `gorm:"foreignKey:UserID" json:"-"`
	FeedbackID     uint
	Feedback       Feedback `gorm:"foreignKey:FeedbackID" json:"-"`
	QuestionAnswer []QuestionAnswer
}
type QuestionAnswer struct {
	gorm.Model
	QuestionID string
	Answer     string
}
