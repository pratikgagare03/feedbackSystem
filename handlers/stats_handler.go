package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/helper"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"gorm.io/gorm"
)

func GetFeedbackStats(c *gin.Context) {
	// Get feedback id from path parameter
	feedbackID := c.Param("feedbackId")
	if feedbackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Feedback ID is required"})
		return
	}
	if matched, err := helper.MatchFeedbackOwner(c, feedbackID); err != nil {
		logger.Logs.Error().Msgf("Error while matching feedback owner: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !matched {
		logger.Logs.Error().Msg("Unauthorized to access this resource")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to access this resource"})
		return
	}
	// Get feedback from database
	responses, err := repository.GetResponseRepository().FindResponseByFeedbackId(feedbackID)
	if len(responses) == 0 || err == gorm.ErrRecordNotFound {
		logger.Logs.Error().Msgf("no responses found for feedback Error:%v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "no responses found for feedback"})
		return
	} else if err != nil {
		logger.Logs.Error().Msgf("error while getting responses for feedback: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var qIdToStats = make(map[uint]models.Stats)
	var totalResponses = make(map[uint]int)
	for _, response := range responses {
		totalResponses[response.QuestionID]++
	}

	for _, response := range responses {
		switch response.QuestionType {
		case models.MCQ, models.SingleChoice:
			currStats, exists := qIdToStats[response.QuestionID]
			if !exists {
				options, err := repository.GetOptionsRepository().FindOptionsByQueId(response.QuestionID)
				if err != nil {
					logger.Logs.Error().Msgf("error while getting options count for question: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				var optionsArr []string
				err = json.Unmarshal(options.Options, &optionsArr)
				if err != nil {
					logger.Logs.Error().Msgf("error while unmarshalling options: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				qIdToStats[response.QuestionID] = models.Stats{
					Question:     response.QuestionContent,
					QuestionType: response.QuestionType,
					Stats: models.MCQStats{
						OptionsCount:      make(map[string]int),
						OptionsPercentage: make(map[string]float64),
						AllOptions:        optionsArr,
					},
				}
				currStats = qIdToStats[response.QuestionID]
				mcqStats, _ := currStats.Stats.(models.MCQStats)
				for _, option := range optionsArr {
					mcqStats.OptionsCount[option] = 0
					mcqStats.OptionsPercentage[option] = 0
				}
				qIdToStats[response.QuestionID] = models.Stats{
					Question:     currStats.Question,
					QuestionType: currStats.QuestionType,
					Stats:        mcqStats,
				}
			}

			mcqStats, ok := currStats.Stats.(models.MCQStats)
			mcqStats.TotalResponses = totalResponses[response.QuestionID]
			if !ok {
				logger.Logs.Error().Msg("ERROR: type assertion to MCQStats failed")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}

			optionsArr := strings.Split(response.Answer, os.Getenv("MCQ_DELIMMITER"))
			for _, option := range optionsArr {
				mcqStats.OptionsCount[option]++
				mcqStats.OptionsPercentage[option] = helper.CalculatePercentage(mcqStats.OptionsCount[option], totalResponses[response.QuestionID])
			}
			qIdToStats[response.QuestionID] = models.Stats{
				Question:     currStats.Question,
				QuestionType: currStats.QuestionType,
				Stats:        mcqStats,
			}

		case models.Ratings:
			currStats, exists := qIdToStats[response.QuestionID]
			if !exists {
				var rRange models.RatingsRange
				res := repository.Db.Find(&rRange, "que_id =?", response.QuestionID)
				if res.Error != nil {
					if res.Error == gorm.ErrRecordNotFound {
						logger.Logs.Error().Msg("ERROR: ratings range not found")
						c.JSON(http.StatusNotFound, gin.H{"error": "ratings range not found"})
						return
					}
					logger.Logs.Error().Msg("ERROR: failed to get ratings from db")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ratings from db"})
					return
				}
				qIdToStats[response.QuestionID] = models.Stats{
					Question:     response.QuestionContent,
					QuestionType: response.QuestionType,
					Stats: models.RatingsStats{
						MaxRatingsRange:   rRange.MaxRatingsRange,
						RatingsCount:      make(map[string]int),
						RatingsPercentage: make(map[string]float64),
						TotalResponses:    totalResponses[response.QuestionID],
						AverageRating:     0,
					},
				}
				currStats = qIdToStats[response.QuestionID]
			}

			ratingsStats, ok := currStats.Stats.(models.RatingsStats)
			if !ok {
				logger.Logs.Error().Msg("ERROR: type assertion to RatingsStats failed")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				return
			}

			ratingsStats.RatingsCount[response.Answer]++
			ratingsStats.RatingsPercentage[response.Answer] = helper.CalculatePercentage(
				ratingsStats.RatingsCount[response.Answer],
				totalResponses[response.QuestionID],
			)

			avg := helper.CalculateAverage(ratingsStats.RatingsCount)
			ratingsStats.AverageRating = avg

			qIdToStats[response.QuestionID] = models.Stats{
				Question:     currStats.Question,
				QuestionType: currStats.QuestionType,
				Stats:        ratingsStats,
			}
		case models.TextInput:
			currStats, exists := qIdToStats[response.QuestionID]
			if !exists {
				qIdToStats[response.QuestionID] = models.Stats{
					Question:     response.QuestionContent,
					QuestionType: response.QuestionType,
					Stats: models.TextInputStats{

						TotalResponses: totalResponses[response.QuestionID],
					},
				}
				currStats = qIdToStats[response.QuestionID]
			}
			textInputStats := currStats.Stats.(models.TextInputStats)
			textInputStats.Answers = append(textInputStats.Answers, response.Answer)
			textInputStats.TotalResponses = totalResponses[response.QuestionID]
			qIdToStats[response.QuestionID] = models.Stats{
				Question:     currStats.Question,
				QuestionType: currStats.QuestionType,
				Stats:        textInputStats,
			}
		}
	}

	c.JSON(http.StatusOK, qIdToStats)
}
