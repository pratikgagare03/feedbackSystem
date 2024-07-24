package models

import (
	"gorm.io/gorm"
)

type Feedback struct {
	gorm.Model
	UserID    uint
	Questions []Question
	User      User `gorm:"foreignKey:UserID"`
}
