package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/helper"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"github.com/pratikgagare03/feedback/utils"
)

func CreateFeedback(c *gin.Context) {
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

	finalFeedback.UserID = c.GetUint("uid")

	err = repository.GetFeedbackRepository().InsertFeedback(&finalFeedback)
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
		question.QuestionContent = questionInput.QuestionContent
		question.QuestionType = qtype
		repository.GetQuestionRepository().InsertQuestion(&question)

		switch qtype {
		case models.MCQ:
			{
				var options models.Options
				if len(questionInput.Options) < 2 {
					log.Println("ERROR: atleast 2 options required in mcq")
					c.JSON(http.StatusBadRequest, gin.H{"error": "atleast 2 options required in mcq"})
					return
				}
				options.QueId = question.ID
				options.Options, err = json.Marshal(questionInput.Options)
				if err != nil {
					log.Println("error while marshaling options to []byte")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "error while marshaling options to []byte"})
					return
				}
				repository.GetOptionsRepository().InsertOptions(&options)
			}
		case models.Ratings:
			{
				var rRange models.RatingsRange
				rRange.QueId = question.ID
				if questionInput.MaxRatingsRange < 0 {
					rRange.MaxRatingsRange = 5
				} else {
					rRange.MaxRatingsRange = questionInput.MaxRatingsRange
				}
				res := repository.Db.Create(rRange)
				if res.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save Ratings"})
					return
				}
			}

		}

	}
	c.JSON(http.StatusCreated, finalFeedback)
}

func GetFeedback(c *gin.Context) {
	feedbackId := c.Param("feedbackId")
	if ok, err := utils.IsValidFeedbackId(feedbackId); !ok {
		log.Printf("ERROR:invalid feedbackId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	feedback, err := repository.GetQuestionRepository().FindQuestionsDetailed(feedbackId)
	if err != nil || len(feedback) == 0 {
		log.Printf("ERROR %+v, foundFeedback length:%+v", err, len(feedback))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusFound, feedback)
}
