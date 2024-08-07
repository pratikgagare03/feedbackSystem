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
	"gorm.io/gorm"
)

func SaveFeedbackResponse(c *gin.Context) {
	logger.Logs.Info().Msg("Saving feedback response")
	feedbackID := c.Param("feedbackId")
	userId := c.GetUint("uid")
	// check if the feedbackId is valid as well as the present in the db
	if ok, err := helper.IsValidFeedbackId(feedbackID); ok {
		logger.Logs.Info().Msg("FeedbackId is valid")
	} else if err != nil {
		logger.Logs.Error().Msgf("ERROR:invalid feedbackId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check if the user is valid
	if exists, err := helper.ResponseExistForUser(feedbackID, userId); exists {
		logger.Logs.Error().Msg("ERROR: A Response already exist.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ERROR: A Response already exist."})
		return
	} else if err != nil {
		logger.Logs.Error().Msg("ERROR: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if ok, err := helper.IsFeedbackPublished(feedbackID); err != nil {
		logger.Logs.Error().Msgf("Error while checking feedback published: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		logger.Logs.Error().Msg("Feedback is not published")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feedback is not published"})
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
		responseDb.QuestionID = qna.QuestionID
		// check if the question is present in the feedback
		if question, err := repository.GetQuestionRepository().FindQuestionByQuestionIdFeedbackId(qna.QuestionID, feedbackID); err != nil {
			logger.Logs.Error().Msg("Question with provided id not present in respective feedback.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question with provided id not present in respective feedback."})
			return
		} else {
			switch question.QuestionType {
			case models.MCQ, models.SingleChoice:
				{ // check if the answer is present in the options
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
					//split the answer by delimiter for multiple choice
					answerArr := strings.Split(qna.Answer, os.Getenv("MCQ_DELIMMITER"))

					for _, answer := range answerArr {
						optionFound := false
						for _, option := range optionsArr {
							if answer == option {
								optionFound = true
								break
							}
						}
						// if option not found return error
						if !optionFound {
							logger.Logs.Error().Msg("ERROR: invalid option selected")
							c.JSON(http.StatusBadRequest, gin.H{"error": "please select a valid option"})
							return
						}
					}
				}
			case models.Ratings:
				{
					// check if the answer is a number
					answerInt, err := strconv.Atoi(qna.Answer)
					if err != nil {
						logger.Logs.Error().Msg("ERROR: answer must be a string(number) for ratings")
						c.JSON(http.StatusBadRequest, gin.H{"error": "answer must be a string(number) for ratings"})
						return
					}
					// check if the answer is in the ratings range
					var rRange models.RatingsRange
					res := repository.Db.Find(&rRange, "que_id =?", question.ID)
					if res.Error != nil {
						logger.Logs.Error().Msg("ERROR: failed to get ratings from db")
						c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ratings from db"})
						return
					}
					// if answer is out of range return error
					if answerInt > rRange.MaxRatingsRange {
						logger.Logs.Error().Msg("ERROR: Rating out of range")
						c.JSON(http.StatusBadRequest, gin.H{"error": "Rating out of range"})
						return
					}
				}
			}
			// save the response
			responseDb.QuestionContent = question.QuestionContent
			responseDb.QuestionType = question.QuestionType
			responseDb.Answer = qna.Answer
			arrResponseDb = append(arrResponseDb, responseDb)
		}
	}
	// save the response array
	err = repository.GetResponseRepository().InsertResponse(arrResponseDb)
	if err != nil {
		logger.Logs.Error().Msg("ERROR: failed to save response " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving response."})
		return
	}
	// return success
	logger.Logs.Info().Msg("Response saved successfully")
	c.JSON(http.StatusCreated, "Your Response has been submitted")
}
func GetAllResponsesForUser(c *gin.Context) {
	logger.Logs.Info().Msg("Fetching Responses for user")
	userId := c.GetUint("uid")

	var responses []models.FeedbackResponse
	var err error
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")

	if dateFrom == "" && dateTo == "" {
		logger.Logs.Info().Msg("dateFrom or dateTo is empty skip applying filter")
		responses, err = repository.GetResponseRepository().GetAllResponsesForUser(userId)
	} else {
		dateFromParsed, dateToParsed, err := helper.GetParsedDateRange(dateFrom, dateTo)
		if err != nil {
			logger.Logs.Error().Msgf("error while parsing date range: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		responses, err = repository.GetResponseRepository().GetAllResponsesForUserDateFilter(userId, dateFromParsed, dateToParsed)
	}

	if len(responses) == 0 || err == gorm.ErrRecordNotFound {
		logger.Logs.Error().Msgf("no responses found for feedback Error:%v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "no responses found for feedback"})
		return
	} else if err != nil {
		logger.Logs.Error().Msgf("error while getting responses for feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var responseOp models.FeedbackResponseOutputForUser
	responseOp.UserID = responses[0].UserID
	for i := 0; i < len(responses); i++ {
		var newResponseOp models.QuestionAnswerForUser
		newResponseOp.FeedbackID = responses[i].FeedbackID
		newResponseOp.CreatedAt = responses[i].CreatedAt
		newResponseOp.UpdatedAt = responses[i].UpdatedAt
		var j int
		for j = i; j < len(responses) && responses[i].FeedbackID == responses[j].FeedbackID; j++ {
			newResponseOp.QnA = append(newResponseOp.QnA, models.QuestionAnswerWithQuestion{Question: responses[j].QuestionContent, Answer: responses[j].Answer})
		}
		if j != i {
			i = j - 1
		}
		responseOp.Responses = append(responseOp.Responses, newResponseOp)
	}
	responseOp.TotalResponses = len(responseOp.Responses)
	c.JSON(http.StatusFound, responseOp)
}

func GetAllResponsesForFeedback(c *gin.Context) {
	logger.Logs.Info().Msg("Fetching Responses for feedback")
	feedbackId := c.Param("feedbackId")
	if feedbackId == "" {
		logger.Logs.Error().Msg("feedback id is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "feedback id is empty"})
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

	if ok, err := helper.IsValidFeedbackId(feedbackId); err != nil {
		logger.Logs.Error().Msgf("Error while checking feedbackId: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if ok {
		logger.Logs.Info().Msg("FeedbackId is valid")
	}
	dateFrom := c.Query("dateFrom")
	dateTo := c.Query("dateTo")

	var responses []models.FeedbackResponse
	var err error
	if dateFrom == "" && dateTo == "" {
		logger.Logs.Info().Msg("dateFrom or dateTo is empty skip applying filter")
		responses, err = repository.GetResponseRepository().FindResponseByFeedbackId(feedbackId)
	} else {
		dateFromParsed, dateToParsed, err := helper.GetParsedDateRange(dateFrom, dateTo)
		if err != nil {
			logger.Logs.Error().Msgf("error while parsing date range: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		responses, err = repository.GetResponseRepository().FindResponseByFeedbackIdDateFilter(feedbackId, dateFromParsed, dateToParsed)
	}

	if len(responses) == 0 || err == gorm.ErrRecordNotFound {
		logger.Logs.Error().Msgf("no responses found for feedback Error:%v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "no responses found for feedback"})
		return
	} else if err != nil {
		logger.Logs.Error().Msgf("error while getting responses for feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var responseOp models.FeedbackResponseOutputForFeedback
	responseOp.FeedbackID = responses[0].FeedbackID

	for i := 0; i < len(responses); i++ {
		var newResponseOp models.QuestionAnswerForFeedback
		newResponseOp.UserID = responses[i].UserID
		newResponseOp.CreatedAt = responses[i].CreatedAt
		newResponseOp.UpdatedAt = responses[i].UpdatedAt
		var j int
		for j = i; j < len(responses) && responses[i].UserID == responses[j].UserID; j++ {
			newResponseOp.QnA = append(newResponseOp.QnA, models.QuestionAnswerWithQuestion{Question: responses[j].QuestionContent, Answer: responses[j].Answer})
		}
		if j != i {
			i = j - 1
		}
		responseOp.Responses = append(responseOp.Responses, newResponseOp)
	}
	responseOp.TotalResponses = len(responseOp.Responses)
	c.JSON(http.StatusFound, responseOp)
}
