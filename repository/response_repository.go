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
	GetAllResponsesForUser(userId uint) ([]models.FeedbackResponse, error)
	GetResponseCountForUser(userId string) (int64, error)
	DeleteResponse(responseID string) error
}

type postgresResponseRepository struct {
	postgresDb *gorm.DB
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
