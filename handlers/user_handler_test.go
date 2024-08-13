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

var testUser models.User

func TestSignUp(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DBNAME", "test")
	router := setupRouter()
	logger.StartLogger()
	logger.Logs.Info().Msg("starting logger for tests")
	repository.Connect()
	t.Run("CreateTestData", func(t *testing.T) {
		user := models.User{
			Email: "test1@example.com",

			Password:   "password123",
			First_name: "Test",
			User_type:  "USER",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		json.Unmarshal(w.Body.Bytes(), &user)
		testUser = user
	})
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
		w := httptest.NewRecorder()
		err := repository.GetUserRepository().DeleteUserByEmail(user.Email)
		router.ServeHTTP(w, req)
		if err != nil {
			logger.Logs.Error().Msgf("Error while deleting user: %v", err)
		}
		t.Logf("Response Status: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "user created successfully")
	})

	t.Run("TestSignUpFailMissingRequiredFields", func(t *testing.T) {
		user := models.User{
			Email: "test2@example.com",
		}

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
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

func AuthMiddlewareSet(userType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("email", "test1@example.com")
		c.Set("user_type", userType)
		c.Set("uid", testUser.ID)
		c.Next()
	}
}
func TestGetUsers(t *testing.T) {
	t.Run("TestGetUsers_Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		r := gin.Default()
		r.Use(AuthMiddlewareSet("ADMIN"))
		r.GET("/users", GetUsers)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("TestGetUsersFail_Normal_User", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		// Create a new Gin router and apply the middleware and handler
		r := gin.Default()
		r.Use(AuthMiddlewareSet("USER"))
		r.GET("/users", GetUsers)
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("TestGetUserAdmin_Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(AuthMiddlewareSet("ADMIN"))
		r.GET("/users/:user_id", GetUser)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/31", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("TestGetUserSameUser_Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(AuthMiddlewareSet("USER"))
		r.GET("/users/:user_id", GetUser)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/31", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("TestGetUserFail_Normal_User", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.Use(AuthMiddlewareSet("USER"))
		r.GET("/users/:user_id", GetUser)
		req, err := http.NewRequest("GET", "/users/23", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
