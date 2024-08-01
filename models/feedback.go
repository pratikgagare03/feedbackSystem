package models

import (
	"gorm.io/gorm"
)

type Feedback struct {
	gorm.Model
	UserID uint `gorm:"column:user_id" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`
}
type FeedbackInput struct {
	Questions []QuestionDetailed
}
