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
	Feedback        Feedback `gorm:"foreignKey:FeedbackID" json:"-"`
	QuestionContent []byte
	QuestionType    QuestionType
}
type QuestionInput struct {
	QuestionContent string
	QuestionType    string
	Options         []string
}

type McqQuestionContent struct {
	QuestionContent string
	Options         []string
}

type QuestionDetailed struct {
	QuestionId      uint
	QuestionContent string
	QuestionType    QuestionType
	Options         []string
}
