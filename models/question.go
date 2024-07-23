package models

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	QuestionContent string
	QuestionType    string
	Status          string
}

