package models

import "gorm.io/gorm"

type QuestionType string

const (
	SingleChoice QuestionType = "singlechoice"
	MCQ          QuestionType = "mcq"
	TextInput    QuestionType = "textinput"
	Ratings      QuestionType = "ratings"
)

// Question table in database.
type Question struct {
	gorm.Model
	FeedbackID      uint         `json:"feedback_id"`
	Feedback        Feedback     `gorm:"foreignKey:FeedbackID" json:"-"`
	QuestionContent string       `json:"question_content" validate:"required"`
	QuestionType    QuestionType `validate:"required,eq=mcq|eq=textinput|eq=ratings|eq=singlechoice"`
}

// QuestionDetailed is a struct that represents a detailed question for getFeedback endpoint.
type QuestionDetailed struct {
	QuestionId      uint     `json:"question_id"`
	QuestionContent string   `json:"question_content"`
	QuestionType    string   `validate:"required,eq=mcq|eq=textinput|eq=ratings"`
	Options         []string `json:"options"  omitempty:"true"`
	MaxRatingsRange int      `json:"max_ratings_range" omitempty:"true"`
}

// options table in database
type Options struct {
	QueId    uint     `json:"que_id" validate:"required"`
	Question Question `gorm:"foreignKey:QueId" json:"-" validate:"required"`
	Options  []byte   `json:"options" validate:"required"`
}

// ratings_ranges table in database
type RatingsRange struct {
	QueId           uint     `json:"que_id" validate:"required"`
	Question        Question `gorm:"foreignKey:QueId" json:"-" validate:"required"`
	MaxRatingsRange int      `validate:"required"`
}
