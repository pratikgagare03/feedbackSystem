package models

import (
	"time"

	"gorm.io/gorm"
) // FeedbackResponse table in database.
type FeedbackResponse struct {
	gorm.Model
	UserID          uint         `json:"user_id"`
	User            User         `gorm:"foreignKey:UserID" json:"-"`
	FeedbackID      uint         `json:"feedback_id"`
	Feedback        Feedback     `gorm:"foreignKey:FeedbackID" json:"-"`
	QuestionID      uint         `json:"question_id"`
	Question        Question     `gorm:"foreignKey:QuestionID" json:"-"`
	QuestionContent string       `json:"question_content"`
	QuestionType    QuestionType `json:"question_type"`
	Answer          string       `json:"answer"`
}

// FeedbackResponseInput is a struct that represents the input for feedback response.
type FeedbackResponseInput struct {
	QuestionAnswer []QuestionAnswer `json:"question_answer" validate:"required"`
}
type QuestionAnswer struct {
	QuestionID uint   `json:"question_id" validate:"required"`
	Answer     string `json:"answer" validate:"required"`
}
type QuestionAnswerWithQuestion struct {
	Question string `json:"question"`
	Answer   string `json:"answer" validate:"required"`
}

// QuestionAnswer is a struct that represents the answer for a question.

type FeedbackResponseOutputForUser struct {
	UserID         uint                    `json:"user_id"`
	TotalResponses int                     `json:"total_responses"`
	Responses      []QuestionAnswerForUser `json:"Responses"`
}
type QuestionAnswerForUser struct {
	FeedbackID uint                         `json:"feedback_id"`
	CreatedAt  time.Time                    `json:"created_at"`
	UpdatedAt  time.Time                    `json:"updated_at"`
	QnA        []QuestionAnswerWithQuestion `json:"QnA"`
}

type FeedbackResponseOutputForFeedback struct {
	FeedbackID     uint                        `json:"feedback_id"`
	TotalResponses int                         `json:"total_responses"`
	Responses      []QuestionAnswerForFeedback `json:"Responses"`
}

type QuestionAnswerForFeedback struct {
	UserID    uint                         `json:"user_id"`
	CreatedAt time.Time                    `json:"created_at"`
	UpdatedAt time.Time                    `json:"updated_at"`
	QnA       []QuestionAnswerWithQuestion `json:"QnA"`
}
