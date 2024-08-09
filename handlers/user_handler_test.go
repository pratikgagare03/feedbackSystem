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

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/signup", SignUp)
	router.POST("/login", Login)
	router.GET("/users", GetUsers)
	return router
}

func TestSignUp(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DBNAME", "test")
	router := setupRouter()
	logger.StartLogger()
	repository.Connect()
	t.Run("TestSignUpSuccess", func(t *testing.T) {
		user := models.User{
			Email: "test1@example.com",

			Password:   "password123",
			First_name: "Test",
			User_type:  "USER",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		err := repository.GetUserRepository().DeleteUserByEmail(user.Email)
		if err != nil {
			logger.Logs.Error().Msgf("Error while deleting user: %v", err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("Response Status: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "user created successfully")
	})

	t.Run("TestSignUpFailMissingRequiredFields", func(t *testing.T) {
		user := models.User{
			Email: "test1@example.com",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		err := repository.GetUserRepository().DeleteUserByEmail(user.Email)
		if err != nil {
			logger.Logs.Error().Msgf("Error while deleting user: %v", err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		t.Logf("Response Status: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestSignUpFailUserExists", func(t *testing.T) {
		user := models.User{
			Email:      "existing@example.com",
			Password:   "password123",
			User_type:  "USER",
			First_name: "test",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}

func TestLogin(t *testing.T) {
	router := setupRouter()
	t.Run("TestLoginSuccess", func(t *testing.T) {

		user := models.User{
			Email:    "existing@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("TestLoginFailure_InvalidCredentials", func(t *testing.T) {

		user := models.User{
			Email:    "existing@example.com",
			Password: "wrongpassword",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestLoginFailure_MissingCredentials", func(t *testing.T) {

		user := models.User{
			Email: "existing@example.com",
			// Missing Password
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestLoginFailure_InvalidEmailFormat", func(t *testing.T) {

		user := models.User{
			Email:    "invalid-email",
			Password: "password123",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestLoginFailure_EmptyPassword", func(t *testing.T) {

		user := models.User{
			Email:    "existing@example.com",
			Password: "", // Empty password
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("TestLoginFailure_UserNotFound", func(t *testing.T) {

		user := models.User{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetUsers(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, 1.0, response["totalCount"])
}

// func TestGetUser_Success(t *testing.T) {
// 	router := setupRouter()
// 	router.GET("/users/:user_id", GetUser)

// 	// Mock the repository to return a user
// 	req, _ := http.NewRequest("GET", "/users/1", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	var user models.User
// 	json.Unmarshal(w.Body.Bytes(), &user)
// 	assert.Equal(t, "test@example.com", user.Email)
// }

// func TestGetUser_NotFound(t *testing.T) {
// 	router := setupRouter()
// 	router.GET("/users/:user_id", GetUser)

// 	req, _ := http.NewRequest("GET", "/users/1", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// }
