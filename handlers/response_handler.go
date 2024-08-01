package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"github.com/pratikgagare03/feedback/utils"
	"gorm.io/gorm"
)

func SaveFeedbackResponse(c *gin.Context) {
	feedbackID := c.Param("feedbackId")
	userId := c.GetUint("uid")
	if ok, err := utils.IsValidFeedbackId(feedbackID); !ok {
		logger.Logs.Error().Msgf("ERROR:invalid feedbackId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok := utils.ResponseExistForUser(feedbackID, userId); ok {
		logger.Logs.Error().Msg("ERROR: A Response already exist.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ERROR: A Response already exist."})
		return
	}

	var responseInput models.FeedbackResponseInput
	if err := c.ShouldBindJSON(&responseInput); err != nil {
		logger.Logs.Error().Msgf("ERROR: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := validate.Struct(responseInput)
	if err != nil {
		logger.Logs.Error().Msgf("ERROR: %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedbackIdInt, _ := strconv.Atoi(feedbackID)

	var arrResponseDb []models.FeedbackResponse

	for _, qna := range responseInput.QuestionAnswer {
		var responseDb models.FeedbackResponse
		responseDb.UserID = userId
		responseDb.FeedbackID = uint(feedbackIdInt)
		//working here get question match the options
		if question, err := repository.GetQuestionRepository().FindQuestionByQuestionIdFeedbackId(qna.QuestionID, feedbackID); err == gorm.ErrRecordNotFound {
			log.Printf("ERROR:%+v Question with provided id not present in respective feedback.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question with provided id not present in respective feedback."})
			return
		} else {
			responseDb.QuestionType = question.QuestionType
			responseDb.QuestionContent = question.QuestionContent
			switch question.QuestionType {
			case models.MCQ:
				{
					options, err := repository.GetOptionsRepository().FindOptionsByQueId(qna.QuestionID)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					var optionsArr []string
					err = json.Unmarshal(options.Options, &optionsArr)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					optionFound := false
					for _, option := range optionsArr {
						if option == qna.Answer {
							optionFound = true
							break
						}
					}
					if !optionFound {
						c.JSON(http.StatusBadRequest, gin.H{"error": "please select a valid option"})
						return
					}
				}
			case models.Ratings:
				{
					answerInt, err := strconv.Atoi(qna.Answer)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "answer must be a string(number) for ratings"})
						return
					}
					var rRange models.RatingsRange
					res := repository.Db.Find(&rRange, "que_id =?", question.ID)
					if res.Error != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ratings from db"})
						return
					}

					if answerInt > rRange.MaxRatingsRange {
						log.Printf("answer rating set to max as recieved rating vas higher that maxRating")
						qna.Answer = strconv.Itoa(rRange.MaxRatingsRange)
					}
				}
			}
			responseDb.Answer = qna.Answer
			arrResponseDb = append(arrResponseDb, responseDb)
		}
	}

	err = repository.GetResponseRepository().InsertResponse(arrResponseDb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving response."})
		return
	}

	c.JSON(http.StatusCreated, "Your Response has been submitted")
}

func GetAllResponses(c *gin.Context) {
	userId := c.Param("userID")
	if ok, err := utils.IsValidUser(userId); !ok {
		log.Printf("ERROR:invalid userId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	feedbackID := c.Param("feedbackId")
	if ok, err := utils.IsValidFeedbackId(feedbackID); !ok {
		log.Printf("ERROR:invalid feedbackId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok := utils.ResponseExistForFeedback(feedbackID); !ok {
		log.Print("ERROR: No responses present.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ERROR: No responses present."})
		return
	}
	//working
}
