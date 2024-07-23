package main

import (
	"log"

	"github.com/pratikgagare03/feedback/database"
	"github.com/pratikgagare03/feedback/models"
)

func main() {
	db, err := database.Setup()
	if err != nil {
		log.Println(err)
	}
	db.AutoMigrate(models.Question{})
	q1 := models.Question{
		QuestionContent: "How are you",
		QuestionType:    "mcq",
		Status:          "unattempted",
	}
	db.Create(&q1)
}
