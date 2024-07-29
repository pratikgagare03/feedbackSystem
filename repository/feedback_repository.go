package repository

import (
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
	InsertFeedback(feedback *models.Feedback) error
	FindFeedbackByID(feedbackID string) (models.Feedback, error)
	UpdateFeedback(feedback *models.Feedback) error
	DeleteFeedback(feedbackID string) error
	GetFeedbacks(tagcontains string) ([]models.Feedback, error)
}

type postgresFeedbackRepository struct {
	postgresDb *gorm.DB
}

// DeleteFeedback implements FeedbackRepository.
func (p *postgresFeedbackRepository) DeleteFeedback(feedbackID string) error {
	panic("unimplemented")
}

// FindFeedbackByID implements FeedbackRepository.
func (p *postgresFeedbackRepository) FindFeedbackByID(feedbackID string) (models.Feedback, error) {
	var fd models.Feedback
	res := Db.First(&fd, feedbackID)
	return fd, res.Error
}

// GetFeedbacks implements FeedbackRepository.
func (p *postgresFeedbackRepository) GetFeedbacks(tagcontains string) ([]models.Feedback, error) {
	panic("unimplemented")
}

// InsertFeedback implements FeedbackRepository.
func (p *postgresFeedbackRepository) InsertFeedback(feedback *models.Feedback) error {
	res := Db.Create(&feedback)
	return res.Error
}

// UpdateFeedback implements FeedbackRepository.
func (p *postgresFeedbackRepository) UpdateFeedback(feedback *models.Feedback) error {
	panic("unimplemented")
}

func newPostgresFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &postgresFeedbackRepository{
		postgresDb: db,
	}
}
func GetFeedbackRepository() FeedbackRepository {
	return newPostgresFeedbackRepository(Db)
}
