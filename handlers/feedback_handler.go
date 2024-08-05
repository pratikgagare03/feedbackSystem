package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/helper"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
)

// CreateFeedback creates a new feedback
func CreateFeedback(c *gin.Context) {
	logger.Logs.Info().Msg("Creating feedback")

	// Check if the user is an admin
	err := helper.CheckUserType(c, "ADMIN")
	if err != nil {
		logger.Logs.Error().Msgf("Error while checking user type: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var newFeedback models.FeedbackInput
	// Bind the json to the struct
	if err := c.ShouldBindJSON(&newFeedback); err != nil {
		logger.Logs.Error().Msgf("Error while binding json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the input
	err = validate.Struct(newFeedback)
	if err != nil {
		logger.Logs.Error().Msgf("Error while validating feedback: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the questions are empty
	if len(newFeedback.Questions) == 0 {
		logger.Logs.Error().Msg("atleast 1 questions required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "atleast 1 questions required"})
		return
	}
	var finalFeedback models.Feedback

	// Copy the values from the input to the final feedback
	finalFeedback.UserID = c.GetUint("uid")
	// Insert the feedback
	err = repository.GetFeedbackRepository().InsertFeedback(&finalFeedback)
	if err != nil {
		logger.Logs.Error().Msgf("Error while inserting feedback: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Iterate over the questions and inset them after validating
	for _, questionInput := range newFeedback.Questions {
		// check for empty question
		if len(questionInput.QuestionContent) == 0 {
			logger.Logs.Error().Msgf("Question cannot be empty, QuestionContent:%+v", questionInput.QuestionContent)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question cannot be empty"})
			return
		}

		var question models.Question
		// Get the question type from Models
		qtype, err := helper.GetQuestionType(questionInput.QuestionType)
		if err != nil {
			logger.Logs.Error().Msgf("Error while getting question type: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Copy the values from the input to the question
		question.FeedbackID = finalFeedback.ID
		question.QuestionContent = questionInput.QuestionContent
		question.QuestionType = qtype
		// Insert the question
		err = repository.GetQuestionRepository().InsertQuestion(&question)
		if err != nil {
			logger.Logs.Error().Msgf("Error while inserting question: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert the options or ratings based on the question type
		switch qtype {
		case models.MCQ, models.SingleChoice:
			{
				var options models.Options
				// check for empty options
				if len(questionInput.Options) < 2 {
					logger.Logs.Error().Msg("atleast 2 options required in mcq")
					c.JSON(http.StatusBadRequest, gin.H{"error": "atleast 2 options required in mcq"})
					return
				}
				options.QueId = question.ID
				// Marshal the options to []byte
				options.Options, err = json.Marshal(questionInput.Options)
				if err != nil {
					logger.Logs.Error().Msgf("Error while marshaling options to []byte: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "error while marshaling options to []byte"})
					return
				}
				// Insert the options in the database
				err = repository.GetOptionsRepository().InsertOptions(&options)
				if err != nil {
					logger.Logs.Error().Msgf("Error while saving options: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save options"})
					return
				}
			}
		case models.Ratings:
			{
				var rRange models.RatingsRange
				rRange.QueId = question.ID
				// check for negative ratings range
				if questionInput.MaxRatingsRange < 0 {
					// if max ratings range is negative, set it to 5
					rRange.MaxRatingsRange = 5
				} else {
					// else set it to the input value
					rRange.MaxRatingsRange = questionInput.MaxRatingsRange
				}
				// Insert the ratings range
				res := repository.Db.Create(rRange)
				if res.Error != nil {
					logger.Logs.Error().Msgf("Error while saving Ratings: %v", res.Error)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save Ratings"})
					return
				}
			}

		}

	}
	// Return the final feedback
	logger.Logs.Info().Msg("Feedback created successfully")
	c.JSON(http.StatusCreated, finalFeedback)
}

func GetFeedback(c *gin.Context) {
	logger.Logs.Info().Msg("Fetching feedback")
	feedbackId := c.Param("feedbackId")
	// Check if the feedbackId is valid i.e present in the database
	if ok, err := helper.IsValidFeedbackId(feedbackId); ok {
		logger.Logs.Info().Msg("FeedbackId is valid")
	} else if err != nil {
		logger.Logs.Error().Msgf("ERROR:invalid feedbackId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Fetch the feedback Questions
	feedback, err := repository.GetQuestionRepository().FindQuestionsDetailed(feedbackId)
	if err != nil || len(feedback) == 0 {
		logger.Logs.Error().Msgf("Error while fetching feedback: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Return the feedback
	logger.Logs.Info().Msg("Feedback fetched successfully")
	c.JSON(http.StatusFound, feedback)
}

func TogglePublishStatus(published bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		feedbackId := c.Param("feedbackId")
		if ok, err := helper.IsValidFeedbackId(feedbackId); ok {
			logger.Logs.Info().Msg("FeedbackId is valid")
		} else if err != nil {
			logger.Logs.Error().Msgf("ERROR:invalid feedbackId %+v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if ok, err := helper.MatchFeedbackOwner(c, feedbackId); err != nil {
			logger.Logs.Error().Msgf("Error while matching feedback owner: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if !ok {
			logger.Logs.Error().Msg("Unauthorized to access this resource")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to access this resource"})
			return
		}

		err := repository.GetFeedbackRepository().UpdatePublishedStatus(feedbackId, published)
		if err != nil {
			logger.Logs.Error().Msgf("Error while publishing feedback: %v", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if published {
			logger.Logs.Info().Msg("Feedback published successfully")
			c.JSON(http.StatusOK, gin.H{"message": "Feedback published successfully"})
		} else {
			logger.Logs.Info().Msg("Feedback unpublished successfully")
			c.JSON(http.StatusOK, gin.H{"message": "Feedback unpublished successfully"})
		}
	}
}
