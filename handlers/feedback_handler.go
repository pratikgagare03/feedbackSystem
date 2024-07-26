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
	if ok, err := utils.IsValidUser(userId); !ok {
		log.Printf("ERROR:invalid userId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	if len(newFeedback.Questions) == 0 {
		log.Println("ERROR: atleast 1 questions required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "atleast 1 questions required"})
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
		if len(questionInput.QuestionContent) == 0 {
			log.Println("ERROR: Question cannot be empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question cannot be empty"})
			return
		}

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
				if len(questionInput.Options) < 2 {
					log.Println("ERROR: atleast 2 options required in mcq")
					c.JSON(http.StatusBadRequest, gin.H{"error": "atleast 2 options required in mcq"})
					return
				}

				question.QuestionContent, _ = json.Marshal(models.McqQuestionContent{
					QuestionContent: questionInput.QuestionContent,
					Options:         questionInput.Options,
				})
			}
		default:
			{
				question.QuestionContent, _ = json.Marshal(questionInput.QuestionContent)
			}

		}
		repository.GetQuestionRepository().InsertQuestion(context.TODO(), &question)
	}
	c.JSON(http.StatusCreated, finalFeedback)
}

func GetFeedback(c *gin.Context) {
	userId := c.Param("userId")
	feedbackId := c.Param("feedbackId")
	if ok, err := utils.IsValidUser(userId); !ok {
		log.Printf("ERROR:invalid userId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok, err := utils.IsValidFeedbackId(feedbackId); !ok {
		log.Printf("ERROR:invalid feedbackId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	feedback, err := repository.GetQuestionRepository().GetQuestionsByFeedbackID(context.TODO(), feedbackId)
	if err != nil || len(feedback) == 0 {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusFound, feedback)
}
