package models

import (
	"gorm.io/gorm"
)

type Response struct {
	gorm.Model
	FeedbackID uint
	QuestionID uint
	Response   string
}
