package repository

import (
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type OptionsRepository interface {
	InsertOptions(option *models.Options) error
	FindOptionsByQueId(questionID uint) (*models.Options, error)
	UpdateOptions(option *models.Options) error
	DeleteOptions(optionID string) error
}

type postgresOptionsRepository struct {
	postgresDb *gorm.DB
}

// FindOptionsByQueId implements OptionsRepository.
func (p *postgresOptionsRepository) FindOptionsByQueId(questionID uint) (*models.Options, error) {
	var options models.Options
	res := Db.Where("que_id = ?", questionID).Find(&options)
	return &options, res.Error
}

// DeleteOptions implements OptionsRepository.
func (p *postgresOptionsRepository) DeleteOptions(optionID string) error {
	panic("unimplemented")
}

// InsertOptions implements OptionsRepository.
func (p *postgresOptionsRepository) InsertOptions(option *models.Options) error {
	res := Db.Create(&option)
	return res.Error
}

// UpdateOptions implements OptionsRepository.
func (p *postgresOptionsRepository) UpdateOptions(option *models.Options) error {
	panic("unimplemented")
}

func newPostgresOptionsRepository(db *gorm.DB) OptionsRepository {
	return &postgresOptionsRepository{
		postgresDb: db,
	}
}
func GetOptionsRepository() OptionsRepository {
	return newPostgresOptionsRepository(Db)
}
