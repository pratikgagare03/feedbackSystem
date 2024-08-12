package handlers

import (
	"bytes"
	"encoding/json"
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
	router.POST("/feedback/:feedbackId/publish", TogglePublishStatus(true))
	router.POST("/feedback/:feedbackId/unpublish", TogglePublishStatus(false))
	return router
}

func TestCreateFeedback(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DBNAME", "test")
	gin.SetMode(gin.TestMode)
	logger.StartLogger()
	repository.Connect()

	t.Run("TestCreateFeedbackSuccess", func(t *testing.T) {
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
