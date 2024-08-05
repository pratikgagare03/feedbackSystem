package repository

import (
	"errors"

	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
	InsertFeedback(feedback *models.Feedback) error
	FindFeedbackByID(feedbackID string) (models.Feedback, error)
	UpdatePublishedStatus(feedbackID string, Published bool) error
}

type postgresFeedbackRepository struct {
	postgresDb *gorm.DB
}

// UnpublishFeedback implements FeedbackRepository.
func (p *postgresFeedbackRepository) UpdatePublishedStatus(feedbackID string, Published bool) error {
	var feedback models.Feedback
	res := Db.Find(&feedback, "id=?", feedbackID)
	if res.Error == nil && feedback.Published == Published {
		return errors.New("feedback already in the desired state")
	}
	res = Db.Model(&models.Feedback{}).Where("id=?", feedbackID).Update("published", Published)
	return res.Error
}

// FindFeedbackByID implements FeedbackRepository.
func (p *postgresFeedbackRepository) FindFeedbackByID(feedbackID string) (models.Feedback, error) {
	var fd models.Feedback
	res := Db.First(&fd, "id=?", feedbackID)
	return fd, res.Error
}

// InsertFeedback implements FeedbackRepository.
func (p *postgresFeedbackRepository) InsertFeedback(feedback *models.Feedback) error {
	res := Db.Create(&feedback)
	return res.Error
}

func newPostgresFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &postgresFeedbackRepository{
		postgresDb: db,
	}
}
func GetFeedbackRepository() FeedbackRepository {
	return newPostgresFeedbackRepository(Db)
}
