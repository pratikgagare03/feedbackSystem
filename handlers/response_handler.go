package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/helper"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"github.com/pratikgagare03/feedback/utils"
	"gorm.io/gorm"
)

func SaveFeedbackResponse(c *gin.Context) {
	logger.Logs.Info().Msg("Saving feedback response")
	feedbackID := c.Param("feedbackId")
	userId := c.GetUint("uid")
	if ok, err := utils.IsValidFeedbackId(feedbackID); ok {
		logger.Logs.Info().Msg("FeedbackId is valid")
	} else if err != nil {
		logger.Logs.Error().Msgf("ERROR:invalid feedbackId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if exists, err := utils.ResponseExistForUser(feedbackID, userId); exists {
		logger.Logs.Error().Msg("ERROR: A Response already exist.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ERROR: A Response already exist."})
		return
	} else if err != nil {
		logger.Logs.Error().Msg("ERROR: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	// iterate over the questions and save the response
	for _, qna := range responseInput.QuestionAnswer {
		var responseDb models.FeedbackResponse
		responseDb.UserID = userId
		responseDb.FeedbackID = uint(feedbackIdInt)

		if question, err := repository.GetQuestionRepository().FindQuestionByQuestionIdFeedbackId(qna.QuestionID, feedbackID); err != nil {
			logger.Logs.Error().Msg("Question with provided id not present in respective feedback.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question with provided id not present in respective feedback."})
			return
		} else {
			switch question.QuestionType {
			case models.MCQ, models.SingleChoice:
				{
					options, err := repository.GetOptionsRepository().FindOptionsByQueId(qna.QuestionID)
					if err != nil {
						logger.Logs.Error().Msg("ERROR: failed to get options from db")
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					var optionsArr []string
					err = json.Unmarshal(options.Options, &optionsArr)
					if err != nil {
						logger.Logs.Error().Msg("ERROR: failed to unmarshal options")
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					answerArr := strings.Split(qna.Answer, os.Getenv("MCQ_DELIMMITER"))

					for _, answer := range answerArr {
						optionFound := false
						for _, option := range optionsArr {
							if answer == option {
								optionFound = true
								break
							}
						}
						if !optionFound {
							logger.Logs.Error().Msg("ERROR: invalid option selected")
							c.JSON(http.StatusBadRequest, gin.H{"error": "please select a valid option"})
							return
						}
					}
				}
			case models.Ratings:
				{
					answerInt, err := strconv.Atoi(qna.Answer)
					if err != nil {
						logger.Logs.Error().Msg("ERROR: answer must be a string(number) for ratings")
						c.JSON(http.StatusBadRequest, gin.H{"error": "answer must be a string(number) for ratings"})
						return
					}
					var rRange models.RatingsRange
					res := repository.Db.Find(&rRange, "que_id =?", question.ID)
					if res.Error != nil {
						logger.Logs.Error().Msg("ERROR: failed to get ratings from db")
						c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ratings from db"})
						return
					}

					if answerInt > rRange.MaxRatingsRange {
						logger.Logs.Error().Msg("ERROR: Rating out of range")
						logger.Logs.Info().Msgf("Answer ratimg out of range %d set ratings to max rane %d", answerInt, rRange.MaxRatingsRange)
						//set the answer to max range
						qna.Answer = strconv.Itoa(rRange.MaxRatingsRange)
					}
				}
			}
			responseDb.QuestionContent = question.QuestionContent
			responseDb.QuestionType = question.QuestionType
			responseDb.Answer = qna.Answer
			arrResponseDb = append(arrResponseDb, responseDb)
		}
	}
	err = repository.GetResponseRepository().InsertResponse(arrResponseDb)
	if err != nil {
		logger.Logs.Error().Msg("ERROR: failed to save response " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving response."})
		return
	}
	logger.Logs.Info().Msg("Response saved successfully")
	c.JSON(http.StatusCreated, "Your Response has been submitted")
}
func GetAllResponsesForUser(c *gin.Context) {
	userId := c.Param("userID")
	if err := helper.MatchUserTypeToUid(c, userId); err != nil {
		logger.Logs.Error().Msgf("error while matching user type to uid: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	responses, err := repository.GetResponseRepository().GetAllResponsesForUser(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Logs.Error().Msgf("no responses found for user Error:%v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "no responses found for user"})
			return
		}
		logger.Logs.Error().Msgf("error while getting responses for user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var responseOp models.FeedbackResponseOutput
	responseOp.ID = responses[0].ID
	responseOp.UserID = responses[0].UserID
	responseOp.CreatedAt = responses[0].CreatedAt
	responseOp.UpdatedAt = responses[0].UpdatedAt
	responseOp.DeletedAt = responses[0].DeletedAt

	for i := 0; i < len(responses); i++ {
		var newResponseOp models.QuestionAnswerWithFeedback
		newResponseOp.FeedbackId = responses[i].FeedbackID
		var j int
		for j = i; j < len(responses) && responses[i].FeedbackID == responses[j].FeedbackID; j++ {
			newResponseOp.QnA = append(newResponseOp.QnA, models.QuestionAnswerWithQuestion{Question: responses[j].QuestionContent, Answer: responses[j].Answer})
		}
		if j != i {
			i = j - 1
		}
		responseOp.Responses = append(responseOp.Responses, newResponseOp)
	}
	c.JSON(http.StatusFound, responseOp)
}
