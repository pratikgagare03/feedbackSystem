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
	QuestionContent string
	QuestionType    QuestionType `validate:"required,eq=mcq|eq=textinput|eq=ratings"`
}
type QuestionDetailed struct {
	QuestionId      uint
	QuestionContent string
	QuestionType    string `validate:"required,eq=mcq|eq=textinput|eq=ratings"`
	Options         []string
	MaxRatingsRange int
}
type Options struct {
	QueId    uint
	Question Question `gorm:"foreignKey:QueId" json:"-"`
	Options  []byte   `validate:"required"`
}

type RatingsRange struct {
	QueId           uint
	Question        Question `gorm:"foreignKey:QueId" json:"-"`
	MaxRatingsRange int      `validate:"required"`
}
