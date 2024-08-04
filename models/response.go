package models

import (
	"gorm.io/gorm"
) // FeedbackResponse table in database.
type FeedbackResponse struct {
	gorm.Model
	UserID          uint         `json:"user_id"`
	User            User         `gorm:"foreignKey:UserID" json:"-"`
	FeedbackID      uint         `json:"feedback_id"`
	Feedback        Feedback     `gorm:"foreignKey:FeedbackID" json:"-"`
	QuestionContent string       `json:"question_content"`
	QuestionType    QuestionType `json:"question_type"`
	Answer          string       `json:"answer"`
}

// FeedbackResponseInput is a struct that represents the input for feedback response.
type FeedbackResponseInput struct {
	QuestionAnswer []QuestionAnswer `json:"question_answer" validate:"required"`
}

// QuestionAnswer is a struct that represents the answer for a question.
type QuestionAnswer struct {
	QuestionID uint   `json:"question_id" validate:"required"`
	Answer     string `json:"answer" validate:"required"`
}

type FeedbackResponseOutput struct {
	gorm.Model
	UserID     uint                         `json:"user_id"`
	Responses  []QuestionAnswerWithFeedback `json:"Responses"`
}
type QuestionAnswerWithFeedback struct {
	FeedbackId uint `json:"feedback_id"`
	QnA        []QuestionAnswerWithQuestion
}
type QuestionAnswerWithQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer" validate:"required"`
}
