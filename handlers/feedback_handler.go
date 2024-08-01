package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/helper"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"github.com/pratikgagare03/feedback/utils"
)

func CreateFeedback(c *gin.Context) {
	logger.Logs.Info().Msg("Creating feedback")
	var newFeedback models.FeedbackInput
	if err := c.ShouldBindJSON(&newFeedback); err != nil {
		logger.Logs.Error().Msgf("Error while binding json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := validate.Struct(newFeedback)
	if err != nil {
		logger.Logs.Error().Msgf("Error while validating feedback: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(newFeedback.Questions) == 0 {
		logger.Logs.Error().Msg("atleast 1 questions required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "atleast 1 questions required"})
		return
	}
	var finalFeedback models.Feedback

	finalFeedback.UserID = c.GetUint("uid")

	err = repository.GetFeedbackRepository().InsertFeedback(&finalFeedback)
	if err != nil {
		logger.Logs.Error().Msgf("Error while inserting feedback: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, questionInput := range newFeedback.Questions {
		if len(questionInput.QuestionContent) == 0 {
			logger.Logs.Error().Msg("Question cannot be empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question cannot be empty"})
			return
		}

		var question models.Question
		qtype, err := helper.GetQuestionType(questionInput.QuestionType)
		if err != nil {
			logger.Logs.Error().Msgf("Error while getting question type: %v", err)
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
					logger.Logs.Error().Msg("atleast 2 options required in mcq")
					c.JSON(http.StatusBadRequest, gin.H{"error": "atleast 2 options required in mcq"})
					return
				}
				options.QueId = question.ID
				options.Options, err = json.Marshal(questionInput.Options)
				if err != nil {
					logger.Logs.Error().Msgf("Error while marshaling options to []byte: %v", err)
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
					logger.Logs.Error().Msgf("Error while saving Ratings: %v", res.Error)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save Ratings"})
					return
				}
			}

		}

	}
	logger.Logs.Info().Msg("Feedback created successfully")
	c.JSON(http.StatusCreated, finalFeedback)
}

func GetFeedback(c *gin.Context) {
	logger.Logs.Info().Msg("Fetching feedback")
	feedbackId := c.Param("feedbackId")
	if ok, err := utils.IsValidFeedbackId(feedbackId); !ok {
		logger.Logs.Error().Msgf("Error while validating feedback id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	feedback, err := repository.GetQuestionRepository().FindQuestionsDetailed(feedbackId)
	if err != nil || len(feedback) == 0 {
		logger.Logs.Error().Msgf("Error while fetching feedback: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	logger.Logs.Info().Msg("Feedback fetched successfully")
	c.JSON(http.StatusFound, feedback)
}
