package models

import (
	"gorm.io/gorm"
)

type Feedback struct {
	gorm.Model
	UserID uint
	User   User `gorm:"foreignKey:UserID"`
}
type FeedbackInput struct {
	Questions []QuestionInput
}
