package models

import (
	"gorm.io/gorm"
)

type Feedback struct {
	gorm.Model
	UserID uint `gorm:"column:user_id" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-" omitEmpty:"true"`
	Published bool `json:"published"`
}
type FeedbackInput struct {
	Questions []QuestionDetailed `json:"questions" validate:"required,dive"`
}
