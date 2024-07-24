package models

import "gorm.io/gorm"

type QuestionType string

const (
	MCQ       QuestionType = "mcq"
	TextInput QuestionType = "textinput"
	Ratings   QuestionType = "ratings"
)

type Question struct {
	gorm.Model
	FeedbackID      uint
	Feedback        Feedback `gorm:"foreignKey:FeedbackID"`
	QuestionContent []byte
	QuestionType    QuestionType
}
