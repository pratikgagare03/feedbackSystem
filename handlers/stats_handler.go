package handlers

// import (
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/pratikgagare03/feedback/models"
// )

//pratik working
// func GetFeedbackStats(c *gin.Context) {
// 	// Get feedback id from path parameter
// 	feedbackID, err := strconv.Atoi(c.Param("feedback_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback id"})
// 		return
// 	}

// 	// Get feedback from database
// 	var feedback models.Feedback
// 	result := h.DB.Preload("Questions").First(&feedback, feedbackID)
// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feedback"})
// 		return
// 	}

// 	// Get feedback responses from database
// 	var responses []models.FeedbackResponse
// 	result = h.DB.Where("feedback_id = ?", feedbackID).Find(&responses)
// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching feedback responses"})
// 		return
// 	}

// 	// Calculate statistics
// 	stats := calculateStats(feedback.Questions, responses)

// 	// Return stats
// 	c.JSON(http.StatusOK, stats)
// }

// func calculateStats(questions []models.Question, responses []models.FeedbackResponse) interface{} {
// 	var stats interface{}
// 	// Calculate stats based on question type
// 	for _, question := range questions {
// 		switch question.QuestionType {
// 		case models.MCQ:
// 			stats = calculateMCQStats(question, responses)
// 		case models.Ratings:
// 			stats = calculateRatingsStats(question, responses)
// 		case models.TextInput:
// 			stats = calculateTextInputStats(question, responses)
// 		}
// 	}
// 	return stats
// }

// func calculateMCQStats(question models.Question, responses []models.FeedbackResponse) models.MCQStats {
// 	var mcqStats models.MCQStats
// 	mcqStats.Question = question.QuestionContent
// 	mcqStats.QuestionType = string(question.QuestionType)
// 	mcqStats.OptionsCount = make(map[string]int)
// 	mcqStats.OptionsPercentage = make(map[string]float64)

// 	// Count options
// 	for _, response := range responses {
// 		if response.QuestionContent == question.QuestionContent {
// 			mcqStats.OptionsCount[response.Answer]++
// 		}
// 	}

// 	// Calculate percentage
// 	totalResponses := len(responses)
// 	for option, count := range mcqStats.OptionsCount {
// 		mcqStats.OptionsPercentage[option] = float64(count) / float64(totalResponses) * 100
// 	}

// 	return mcqStats
// }
