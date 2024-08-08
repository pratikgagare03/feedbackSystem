package repository

import (
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type ResponseRepository interface {
	InsertResponse(response []models.FeedbackResponse) error
	FindResponseByID(responseID string) (models.FeedbackResponse, error)
	FindResponseByUserIdFeedbackId(userID uint, feedbackID string) (models.FeedbackResponse, error)
	FindResponseByFeedbackId(feedbackID string) ([]models.FeedbackResponse, error)
	FindResponseByFeedbackIdDateFilter(feedbackID, fromTime, toTime string) ([]models.FeedbackResponse, error)
	GetResponsesWithQueryFilter(params models.QueryParams, feedbackId string) ([]models.FeedbackResponse, error)
	GetAllResponsesForUser(userId uint) ([]models.FeedbackResponse, error)
	GetAllResponsesForUserDateFilter(userId uint, fromTime, toTime string) ([]models.FeedbackResponse, error)
	GetResponseCountForUser(userId string) (int64, error)
	DeleteResponse(responseID string) error
}

type postgresResponseRepository struct {
	postgresDb *gorm.DB
}

func (p *postgresResponseRepository) GetResponsesWithQueryFilter(params models.QueryParams, feedbackId string) ([]models.FeedbackResponse, error) {
	var responses []models.FeedbackResponse
	var err error

	// Construct base query
	query := Db.Model(&models.FeedbackResponse{}).Where("feedback_id = ?", feedbackId)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("DATE(created_at) BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	// Apply text filters
	if params.InQuestion == "true" && params.InAnswer == "true" {
		query = query.Where("question_content ILIKE ? OR answer ILIKE ?", "%"+params.Query+"%", "%"+params.Query+"%")
	} else if params.InQuestion == "true" {
		query = query.Where("question_content ILIKE ?", "%"+params.Query+"%")
	} else if params.InAnswer == "true" {
		query = query.Where("answer ILIKE ?", "%"+params.Query+"%")
	}

	// Execute query
	query = query.Order("created_at ASC").Order("user_id ASC")
	err = query.Find(&responses).Error
	return responses, err
}

// GetResponseCountForUser implements ResponseRepository.
func (p *postgresResponseRepository) GetResponseCountForUser(userId string) (int64, error) {
	var count int64
	res := Db.Model(&models.FeedbackResponse{}).Where("user_id=?", userId).Count(&count)
	return count, res.Error
}

// FindResponseByFeedbackId implements ResponseRepository.
func (p *postgresResponseRepository) FindResponseByFeedbackId(feedbackID string) ([]models.FeedbackResponse, error) {
	var matchingResponses []models.FeedbackResponse
	res := Db.Find(&matchingResponses, "feedback_id=?", feedbackID).Order("user_id")
	return matchingResponses, res.Error
}
func (p *postgresResponseRepository) FindResponseByFeedbackIdDateFilter(feedbackID, fromTime, toTime string) ([]models.FeedbackResponse, error) {
	var matchingResponses []models.FeedbackResponse
	res := Db.Find(&matchingResponses, "feedback_id=? AND DATE(created_at) BETWEEN ? AND ?", feedbackID, fromTime, toTime).Order("user_id")
	return matchingResponses, res.Error
}

// FindResponseByUserIdFeedbackId implements ResponseRepository.
func (p *postgresResponseRepository) FindResponseByUserIdFeedbackId(userID uint, feedbackID string) (models.FeedbackResponse, error) {
	var matchingResponses models.FeedbackResponse
	res := Db.First(&matchingResponses, "feedback_id=? AND user_id=?", feedbackID, userID)
	return matchingResponses, res.Error
}

// DeleteResponse implements ResponseRepository.
func (p *postgresResponseRepository) DeleteResponse(responseID string) error {
	panic("unimplemented")
}

// FindResponseByID implements ResponseRepository.
func (p *postgresResponseRepository) FindResponseByID(responseID string) (models.FeedbackResponse, error) {
	var fd models.FeedbackResponse
	res := Db.First(&fd, responseID)
	return fd, res.Error
}

// GetResponses implements ResponseRepository.
func (p *postgresResponseRepository) GetAllResponsesForUser(userId uint) ([]models.FeedbackResponse, error) {
	var feedbackResponse []models.FeedbackResponse
	res := Db.Find(&feedbackResponse, "user_id=?", userId).Order("feedback_id")
	return feedbackResponse, res.Error
}
func (p *postgresResponseRepository) GetAllResponsesForUserDateFilter(userId uint, fromTime, toTime string) ([]models.FeedbackResponse, error) {
	var feedbackResponse []models.FeedbackResponse
	res := Db.Find(&feedbackResponse, "user_id=? AND DATE(created_at) BETWEEN ? AND ?", userId, fromTime, toTime).Order("feedback_id")
	return feedbackResponse, res.Error
}

// InsertResponse implements ResponseRepository.
func (p *postgresResponseRepository) InsertResponse(response []models.FeedbackResponse) error {
	res := Db.Create(&response)
	return res.Error
}

func newPostgresResponseRepository(db *gorm.DB) ResponseRepository {
	return &postgresResponseRepository{
		postgresDb: db,
	}
}
func GetResponseRepository() ResponseRepository {
	return newPostgresResponseRepository(Db)
}
