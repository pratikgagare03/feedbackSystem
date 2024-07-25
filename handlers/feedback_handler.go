package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/helper"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"github.com/pratikgagare03/feedback/utils"
)

func CreateFeedback(c *gin.Context) {
	userId := c.Param("userId")
	log.Println("hello:", userId)
	if !utils.IsValidUser(userId) {
		log.Printf("ERROR %+v is not a valid userId", userId)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}
	var newFeedback models.FeedbackInput
	if err := c.ShouldBindJSON(&newFeedback); err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := validate.Struct(newFeedback)
	if err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(newFeedback.Questions) < 2 {
		log.Println("ERROR: atleast two questions required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "atleast two questions required"})
		return
	}
	var finalFeedback models.Feedback
	userIdInt, _ := strconv.Atoi(userId)
	finalFeedback.UserID = uint(userIdInt)
	log.Println("hello:", finalFeedback.UserID)

	err = repository.GetFeedbackRepository().InsertFeedback(context.TODO(), &finalFeedback)
	if err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, questionInput := range newFeedback.Questions {
		var question models.Question
		qtype, err := helper.GetQuestionType(questionInput.QuestionType)
		if err != nil {
			log.Printf("ERROR %+v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		question.FeedbackID = finalFeedback.ID
		question.QuestionType = qtype

		switch qtype {
		case models.MCQ:
			{
				question.QuestionContent, _ = json.Marshal(models.McqQusetionContent{
					QuestionContent: questionInput.QuestionContent,
					Options:         questionInput.Options,
				})
			}
		default:
			{
				question.QuestionContent, _ = json.Marshal(questionInput.QuestionContent)
			}

		}

	}
	// repository.GetFeedbackRepository().InsertFeedback(context.TODO(), &newFeedback)
	c.JSON(http.StatusCreated, finalFeedback)
}
