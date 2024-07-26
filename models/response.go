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
	QuestionId string
	Answer     string
}
