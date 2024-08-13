package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"github.com/pratikgagare03/feedback/repository"
	"github.com/stretchr/testify/assert"
)

func setupResponseRouterWithRole(role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(AuthMiddlewareSet(role))
	router.POST("/feedback/:feedbackId/response", SaveFeedbackResponse)
	router.GET("/user/responses", GetAllResponsesForUser)
	router.GET("/feedback/:feedbackId/responses", GetAllResponsesForFeedback)
	return router
}

func TestSaveFeedbackResponse(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DBNAME", "test")
	gin.SetMode(gin.TestMode)
	logger.StartLogger()
	repository.Connect()
	t.Run("TestSaveFeedbackResponseFailInvalidFeedbackId", func(t *testing.T) {
		router := setupResponseRouterWithRole("USER")
		feedbackID := "invalid"
		responseInput := models.FeedbackResponseInput{
			QuestionAnswer: []models.QuestionAnswer{
				{
					QuestionID: 1,          // Set valid question ID
					Answer:     "Option 1", // Or a valid rating
				},
			},
		}

		body, _ := json.Marshal(responseInput)
		req, _ := http.NewRequest("POST", fmt.Sprintf("/feedback/%s/response", feedbackID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestSaveFeedbackResponseFailInvalidData", func(t *testing.T) {
		router := setupResponseRouterWithRole("USER")

		responseInput := models.FeedbackResponseInput{
			QuestionAnswer: []models.QuestionAnswer{
				{
					QuestionID: createdQuestion[0].QuestionId, // Set valid question ID
					Answer:     "",                            // Invalid answer
				},
			},
		}

		body, _ := json.Marshal(responseInput)
		req, _ := http.NewRequest("POST", fmt.Sprintf("/feedback/%v/response", feedbackId), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestSaveFeedbackResponseSuccess", func(t *testing.T) {
		router := setupResponseRouterWithRole("USER")
		responseInput := models.FeedbackResponseInput{
			QuestionAnswer: []models.QuestionAnswer{
				{
					QuestionID: createdQuestion[0].QuestionId, // Set valid question ID
					Answer:     createdQuestion[0].Options[0], // Or a valid rating
				},
			},
		}

		body, _ := json.Marshal(responseInput)
		req, _ := http.NewRequest("POST", fmt.Sprintf("/feedback/%v/response", feedbackId), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `"Your Response has been submitted"`, w.Body.String())
	})
}

func TestGetAllResponsesForUser(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DBNAME", "test")
	gin.SetMode(gin.TestMode)
	logger.StartLogger()
	repository.Connect()

	t.Run("TestGetAllResponsesForUserSuccess", func(t *testing.T) {
		router := setupResponseRouterWithRole("USER")
		userID := 1

		// Mock the user ID in context
		router.Use(func(c *gin.Context) {
			c.Set("uid", userID)
			c.Next()
		})

		req, _ := http.NewRequest("GET", "/user/responses", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("TestGetAllResponsesForUserFailNoResponses", func(t *testing.T) {
		temp := testUser.ID
		testUser.ID = 2
		router := setupResponseRouterWithRole("USER")
		req, _ := http.NewRequest("GET", "/user/responses", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		testUser.ID = temp
	})
}

func TestGetAllResponsesForFeedback(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DBNAME", "test")
	gin.SetMode(gin.TestMode)
	logger.StartLogger()
	repository.Connect()

	t.Run("TestGetAllResponsesForFeedbackSuccess", func(t *testing.T) {
		router := setupResponseRouterWithRole("ADMIN")
		req, _ := http.NewRequest("GET", fmt.Sprintf("/feedback/%v/responses", feedbackId), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("TestGetAllResponsesForFeedbackFailNonOwner", func(t *testing.T) {
		router := setupResponseRouterWithRole("ADMIN")
		feedbackID := "3850"

		req, _ := http.NewRequest("GET", fmt.Sprintf("/feedback/%s/responses", feedbackID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
	t.Run("TestGetAllResponsesForFeedbackNoResponses", func(t *testing.T) {
		router := setupResponseRouterWithRole("ADMIN")
		feedbackID := "3"

		req, _ := http.NewRequest("GET", fmt.Sprintf("/feedback/%s/responses", feedbackID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
