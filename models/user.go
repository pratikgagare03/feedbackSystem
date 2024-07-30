package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId        uuid.UUID `json:"user_id"`
	Username      string    `gorm:"unique"`
	Email         string    `json:"email" validate:"email, required" gorm:"unique"`
	First_name    string    `json:"first_name" validate:"required min=2, max=100"`
	Last_name     string    `json:"last_name"`
	Password      string    `json:"password" validate:"required min=6, max=100"`
	User_type     string    `json:"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Token         string    `json:"token"`
	Refresh_token string    `json:"refresh_token"`
}
