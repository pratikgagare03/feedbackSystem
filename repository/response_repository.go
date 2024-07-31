package repository

import (
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type ResponseRepository interface {
	InsertResponse(response []models.FeedbackResponse) error
	FindResponseByID(responseID string) (models.FeedbackResponse, error)
	FindResponseByUserIdFeedbackId(userID uint, feedbackID string) ([]models.FeedbackResponse, error)
	FindResponseByFeedbackId(feedbackID string) ([]models.FeedbackResponse, error)
	UpdateResponse(response *models.FeedbackResponse) error
	DeleteResponse(responseID string) error
	GetResponses(tagcontains string) ([]models.FeedbackResponse, error)
}

type postgresResponseRepository struct {
	postgresDb *gorm.DB
}

// FindResponseByFeedbackId implements ResponseRepository.
func (p *postgresResponseRepository) FindResponseByFeedbackId(feedbackID string) ([]models.FeedbackResponse, error) {
	var matchingResponses []models.FeedbackResponse
	res := Db.Where("feedback_id = ?", feedbackID).Find(&matchingResponses)
	return matchingResponses, res.Error
}

// FindResponseByUserIdFeedbackId implements ResponseRepository.
func (p *postgresResponseRepository) FindResponseByUserIdFeedbackId(userID uint, feedbackID string) ([]models.FeedbackResponse, error) {
	var matchingResponses []models.FeedbackResponse
	res := Db.Where("feedback_id = ? AND user_id = ?", feedbackID, userID).Find(&matchingResponses)
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
func (p *postgresResponseRepository) GetResponses(tagcontains string) ([]models.FeedbackResponse, error) {
	panic("unimplemented")
}

// InsertResponse implements ResponseRepository.
func (p *postgresResponseRepository) InsertResponse(response []models.FeedbackResponse) error {
	res := Db.Create(&response)
	return res.Error
}

// UpdateResponse implements ResponseRepository.
func (p *postgresResponseRepository) UpdateResponse(response *models.FeedbackResponse) error {
	panic("unimplemented")
}

func newPostgresResponseRepository(db *gorm.DB) ResponseRepository {
	return &postgresResponseRepository{
		postgresDb: db,
	}
}
func GetResponseRepository() ResponseRepository {
	return newPostgresResponseRepository(Db)
}
