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

func setupFeedbackRouterWithRole(role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(AuthMiddlewareSet(role))
	router.POST("/feedback", CreateFeedback)
	router.GET("/feedback/:feedbackId", GetFeedback)
	router.PATCH("/feedback/:feedbackId/publish", TogglePublishStatus(true))
	router.PATCH("/feedback/:feedbackId/unpublish", TogglePublishStatus(false))
	return router
}

var feedbackId uint
var createdQuestion []models.QuestionDetailed

func TestCreateFeedback(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DBNAME", "test")
	gin.SetMode(gin.TestMode)
	logger.StartLogger()
	repository.Connect()

	t.Run("CreateTestData", func(t *testing.T) {
		router := setupRouter()
		user := models.User{
			Email:      "test1@example.com",
			Password:   "password123",
			First_name: "Test",
			User_type:  "ADMIN",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		json.Unmarshal(w.Body.Bytes(), &user)
		user.ID = 1
		testUser = user
	})
	t.Run("TestCreateFeedbackSuccess", func(t *testing.T) {
		logger.Logs.Info().Msg("Test Creating feedback success")
		router := setupFeedbackRouterWithRole("ADMIN")
		feedback := models.FeedbackInput{
			Questions: []models.QuestionDetailed{
				{
					QuestionContent: "What is your favorite color?",
					QuestionType:    "mcq",
					Options:         []string{"Red", "Blue", "Green"},
				},
			},
		}

		body, _ := json.Marshal(feedback)
		req, _ := http.NewRequest("POST", "/feedback", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		var createdFeedback map[string]models.Feedback
		json.Unmarshal(w.Body.Bytes(), &createdFeedback)
		feedbackId = createdFeedback["feedback"].ID
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("TestCreateFeedbackFailNonAdminUser", func(t *testing.T) {
		router := setupFeedbackRouterWithRole("USER")
		feedback := models.FeedbackInput{
			Questions: []models.QuestionDetailed{
				{
					QuestionContent: "What is your favorite color?",
					QuestionType:    "mcq",
					Options:         []string{"Red", "Blue", "Green"},
				},
			},
		}

		body, _ := json.Marshal(feedback)
		req, _ := http.NewRequest("POST", "/feedback", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("TestCreateFeedbackFailInvalidData", func(t *testing.T) {
		router := setupFeedbackRouterWithRole("ADMIN")
		feedback := models.FeedbackInput{
			Questions: []models.QuestionDetailed{
				{
					QuestionContent: "",
					QuestionType:    "mcq",
					Options:         []string{},
				},
			},
		}

		body, _ := json.Marshal(feedback)
		req, _ := http.NewRequest("POST", "/feedback", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestCreateFeedbackFailEmptyOptions", func(t *testing.T) {
		router := setupFeedbackRouterWithRole("ADMIN")
		feedback := models.FeedbackInput{
			Questions: []models.QuestionDetailed{
				{
					QuestionContent: "How are you?",
					QuestionType:    "mcq",
					Options:         []string{},
				},
			},
		}

		body, _ := json.Marshal(feedback)
		req, _ := http.NewRequest("POST", "/feedback", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestCreateFeedbackFailEmptyQuestion", func(t *testing.T) {
		router := setupFeedbackRouterWithRole("ADMIN")
		feedback := models.FeedbackInput{
			Questions: []models.QuestionDetailed{
				{
					QuestionContent: "",
					QuestionType:    "mcq",
					Options:         []string{"fine", "good", "bad"},
				},
			},
		}

		body, _ := json.Marshal(feedback)
		req, _ := http.NewRequest("POST", "/feedback", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("TestCreateFeedbackFailSinleOptionInMcq", func(t *testing.T) {
		router := setupFeedbackRouterWithRole("ADMIN")
		feedback := models.FeedbackInput{
			Questions: []models.QuestionDetailed{
				{
					QuestionContent: "How are you?",
					QuestionType:    "mcq",
					Options:         []string{"fine"},
				},
			},
		}

		body, _ := json.Marshal(feedback)
		req, _ := http.NewRequest("POST", "/feedback", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

func TestGetFeedback_Success(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DBNAME", "test")
	gin.SetMode(gin.TestMode)
	logger.StartLogger()
	repository.Connect()
	t.Run("TestGetFeedbackSuccess", func(t *testing.T) {
		router := setupFeedbackRouterWithRole("ADMIN")
		// Create a request and recorder
		path := fmt.Sprintf("/feedback/%v", feedbackId)
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)
		json.Unmarshal(w.Body.Bytes(), &createdQuestion)
		// Assert the response
		assert.Equal(t, http.StatusFound, w.Code)
	})

	t.Run("TestGetFeedbackFailInvalidFeedbackId", func(t *testing.T) {
		router := setupFeedbackRouterWithRole("ADMIN")
		// Create a request and recorder
		path := fmt.Sprintf("/feedback/%v", 2537)
		req, _ := http.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTogglePublishStatus(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DBNAME", "test")
	gin.SetMode(gin.TestMode)
	logger.StartLogger()
	repository.Connect()
	t.Run("TestTogglePublishStatusSuccess", func(t *testing.T) {
		router := setupFeedbackRouterWithRole("ADMIN")
		// Create a request and recorder
		path := fmt.Sprintf("/feedback/%v/publish", feedbackId)
		req, _ := http.NewRequest(http.MethodPatch, path, nil)
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestTogglePublishStatusFailInvalidFeedbackId", func(t *testing.T) {
		router := setupFeedbackRouterWithRole("ADMIN")
		// Create a request and recorder
		path := fmt.Sprintf("/feedback/%v/publish", "invalid")
		req, _ := http.NewRequest(http.MethodPatch, path, nil)
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)

		// Assert the response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
