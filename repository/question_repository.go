package repository

import (
	"context"

	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	InsertQuestion(ctx context.Context, question *models.Question) error
	FindQuestionByID(ctx context.Context, questionID string) (*models.Question, error)
	UpdateQuestion(ctx context.Context, question *models.Question) error
	DeleteQuestion(ctx context.Context, questionID string) error
	GetQuestions(tagcontains string) ([]models.Question, error)
}

type postgresQuestionRepository struct {
	postgresDb *gorm.DB
}

// DeleteQuestion implements QuestionRepository.
func (p *postgresQuestionRepository) DeleteQuestion(ctx context.Context, questionID string) error {
	panic("unimplemented")
}

// FindQuestionByID implements QuestionRepository.
func (p *postgresQuestionRepository) FindQuestionByID(ctx context.Context, questionID string) (*models.Question, error) {
	panic("unimplemented")
}

// GetQuestions implements QuestionRepository.
func (p *postgresQuestionRepository) GetQuestions(tagcontains string) ([]models.Question, error) {
	panic("unimplemented")
}

// InsertQuestion implements QuestionRepository.
func (p *postgresQuestionRepository) InsertQuestion(ctx context.Context, question *models.Question) error {
	var newQuestion = models.Question{}
	res := Db.Create(&newQuestion)
	return res.Error
}

// UpdateQuestion implements QuestionRepository.
func (p *postgresQuestionRepository) UpdateQuestion(ctx context.Context, question *models.Question) error {
	panic("unimplemented")
}

func newPostgresQuestionRepository(db *gorm.DB) QuestionRepository {
	return &postgresQuestionRepository{
		postgresDb: db,
	}
}
func GetQuestionRepository() QuestionRepository {
	return newPostgresQuestionRepository(Db)
}
