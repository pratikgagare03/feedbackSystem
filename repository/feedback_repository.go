package repository

import (
	"context"

	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
	InsertFeedback(ctx context.Context, feedback *models.Feedback) error
	FindFeedbackByID(ctx context.Context, feedbackID string) (*models.Feedback, error)
	UpdateFeedback(ctx context.Context, feedback *models.Feedback) error
	DeleteFeedback(ctx context.Context, feedbackID string) error
	GetFeedbacks(tagcontains string) ([]models.Feedback, error)
}

type postgresFeedbackRepository struct {
	postgresDb *gorm.DB
}

// DeleteFeedback implements FeedbackRepository.
func (p *postgresFeedbackRepository) DeleteFeedback(ctx context.Context, feedbackID string) error {
	panic("unimplemented")
}

// FindFeedbackByID implements FeedbackRepository.
func (p *postgresFeedbackRepository) FindFeedbackByID(ctx context.Context, feedbackID string) (*models.Feedback, error) {
	panic("unimplemented")
}

// GetFeedbacks implements FeedbackRepository.
func (p *postgresFeedbackRepository) GetFeedbacks(tagcontains string) ([]models.Feedback, error) {
	panic("unimplemented")
}

// InsertFeedback implements FeedbackRepository.
func (p *postgresFeedbackRepository) InsertFeedback(ctx context.Context, feedback *models.Feedback) error {
	res := Db.Create(&feedback)
	return res.Error
}

// UpdateFeedback implements FeedbackRepository.
func (p *postgresFeedbackRepository) UpdateFeedback(ctx context.Context, feedback *models.Feedback) error {
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
