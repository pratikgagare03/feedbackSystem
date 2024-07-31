package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"github.com/pratikgagare03/feedback/utils"
)

func SaveFeedbackResponse(c *gin.Context) {
	feedbackID := c.Param("feedbackId")
	userId := c.GetUint("uid")
	if ok, err := utils.IsValidFeedbackId(feedbackID); !ok {
		log.Printf("ERROR:invalid feedbackId %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ok := utils.ResponseExistForUser(feedbackID, userId); ok {
		log.Print("ERROR: A Response already exist.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "ERROR: A Response already exist."})
		return
	}

	var responseInput models.FeedbackResponseInput
	if err := c.ShouldBindJSON(&responseInput); err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := validate.Struct(responseInput)
	if err != nil {
		log.Printf("ERROR %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	feedbackIdInt, _ := strconv.Atoi(feedbackID)

	var arrResponseDb []models.FeedbackResponse

	for _, qna := range responseInput.QuestionAnswer {
		var responseDb models.FeedbackResponse
		responseDb.UserID = userId
		responseDb.FeedbackID = uint(feedbackIdInt)
		responseDb.QuestionID = qna.QuestionID
		if questions, err := repository.GetQuestionRepository().FindQuestionByQuestionIdFeedbackId(qna.QuestionID, feedbackID); len(questions) == 0 {
			log.Printf("ERROR:%+v Question with provided id not present in respective feedback.", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question with provided id not present in respective feedback."})
			return
		} else {

		}
		responseDb.Answer = qna.Answer
		arrResponseDb = append(arrResponseDb, responseDb)
	}

	repository.GetResponseRepository().InsertResponse(arrResponseDb)
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
