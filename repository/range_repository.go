package repository

import (
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type RatingsRepository interface {
	InsertRatings(feedback *models.RatingsRange) error
	FindRatingsByQueID(feedbackID uint) (models.RatingsRange, error)
}

type postgresRatingsRepository struct {
	postgresdb *gorm.DB
}

// FindRatingsByQueID implements RatingsRepository.
func (p *postgresRatingsRepository) FindRatingsByQueID(que_id uint) (models.RatingsRange, error) {
	var rRange models.RatingsRange
	res := db.Find(&rRange, "que_id =?", que_id)
	return rRange, res.Error
}

// InsertRatings implements RatingsRepository.
func (p *postgresRatingsRepository) InsertRatings(rRange *models.RatingsRange) error {
	res := db.Create(rRange)
	return res.Error
}

func newPostgresRatingsRepository(db *gorm.DB) RatingsRepository {
	return &postgresRatingsRepository{
		postgresdb: db,
	}
}
func GetRatingsRepository() RatingsRepository {
	return newPostgresRatingsRepository(db)
}
