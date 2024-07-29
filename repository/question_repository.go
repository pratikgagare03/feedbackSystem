package repository

import (
	"encoding/json"

	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	InsertQuestion(question *models.Question) error
	FindQuestionByID(questionID string) (*models.Question, error)
	FindQuestionByQuestionIdFeedbackId(questionID string, feedbackID string) ([]models.Question, error)
	GetQuestionsByFeedbackID(questionID string) ([]models.QuestionDetailed, error)
	UpdateQuestion(question *models.Question) error
	DeleteQuestion(questionID string) error
	GetQuestions(tagcontains string) ([]models.Question, error)
}

type postgresQuestionRepository struct {
	postgresDb *gorm.DB
}

// FindQuestionByQuestionIdFeedbackId implements QuestionRepository.
func (p *postgresQuestionRepository) FindQuestionByQuestionIdFeedbackId(questionID string, feedbackID string) ([]models.Question, error) {
	var matchingQuestions []models.Question
	res := Db.Where("id = ? AND feedback_id = ?", questionID, feedbackID).Find(&matchingQuestions)
	return matchingQuestions, res.Error
}

// GetQuestionsByFeedbackID implements QuestionRepository.
func (p *postgresQuestionRepository) GetQuestionsByFeedbackID(feedbackID string) ([]models.QuestionDetailed, error) {
	var questions []models.Question
	res := Db.Where("feedback_id = ?", feedbackID).Find(&questions)

	var quesDetailed []models.QuestionDetailed
	for _, question := range questions {
		var que models.QuestionDetailed
		que.QuestionId = question.ID
		que.QuestionType = question.QuestionType
		if question.QuestionType == models.MCQ {
			var mcqQueContent models.McqQuestionContent
			json.Unmarshal(question.QuestionContent, &mcqQueContent)
			que.QuestionContent = mcqQueContent.QuestionContent
			que.Options = mcqQueContent.Options
		} else {
			var normalQueContent string
			json.Unmarshal(question.QuestionContent, &normalQueContent)
			que.QuestionContent = normalQueContent
		}

		quesDetailed = append(quesDetailed, que)
	}
	return quesDetailed, res.Error
}

// DeleteQuestion implements QuestionRepository.
func (p *postgresQuestionRepository) DeleteQuestion(questionID string) error {
	panic("unimplemented")
}

// FindQuestionByID implements QuestionRepository.
func (p *postgresQuestionRepository) FindQuestionByID(questionID string) (*models.Question, error) {
	panic("unimplemented")
}

// GetQuestions implements QuestionRepository.
func (p *postgresQuestionRepository) GetQuestions(tagcontains string) ([]models.Question, error) {
	panic("unimplemented")
}

// InsertQuestion implements QuestionRepository.
func (p *postgresQuestionRepository) InsertQuestion(question *models.Question) error {
	res := Db.Create(&question)
	return res.Error
}

// UpdateQuestion implements QuestionRepository.
func (p *postgresQuestionRepository) UpdateQuestion(question *models.Question) error {
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
