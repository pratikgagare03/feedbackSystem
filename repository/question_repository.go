package repository

import (
	"encoding/json"

	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	InsertQuestion(question *models.Question) error
	FindQuestionByID(questionID string) (*models.Question, error)
	FindQuestionsDetailed(feedbackId string) ([]models.QuestionDetailed, error)
	FindQuestionByQuestionIdFeedbackId(questionID uint, feedbackID string) (models.Question, error)
	GetQuestionsByFeedbackID(feedbackId string) ([]models.Question, error)
	UpdateQuestion(question *models.Question) error
	DeleteQuestion(questionID string) error
}

type postgresQuestionRepository struct {
	postgresdb *gorm.DB
}

// FindQuestionsDetailed implements QuestionRepository.
func (p *postgresQuestionRepository) FindQuestionsDetailed(feedbackId string) ([]models.QuestionDetailed, error) {
	type questionDetaileddbByte struct {
		ID              uint   `json:"question_id"`
		QuestionContent string `json:"question_content"`
		QuestionType    string `json:"question_type"`
		Options         []byte `json:"options"`
		MaxRatingsRange int    `json:"max_ratings_range"`
	}
	var questionsDetailedByte []questionDetaileddbByte

	res := db.Raw("SELECT q.id,q.question_content,q.question_type, o.options, rr.max_ratings_range FROM questions q LEFT JOIN options o ON q.id = o.que_id LEFT JOIN ratings_ranges rr ON q.id = rr.que_id WHERE (q.feedback_id = ?) AND (q.id = o.que_id OR q.id = rr.que_id OR q.question_type = ?);", feedbackId, models.TextInput).Scan(&questionsDetailedByte)

	var questionsDetailedString []models.QuestionDetailed
	for _, questionDetailedByte := range questionsDetailedByte {
		var q models.QuestionDetailed
		q.QuestionId = questionDetailedByte.ID
		q.QuestionContent = questionDetailedByte.QuestionContent
		q.MaxRatingsRange = questionDetailedByte.MaxRatingsRange
		q.QuestionType = questionDetailedByte.QuestionType
		var optionsString []string
		json.Unmarshal(questionDetailedByte.Options, &optionsString)
		q.Options = optionsString

		questionsDetailedString = append(questionsDetailedString, q)
	}
	return questionsDetailedString, res.Error
}

// DeleteQuestion implements QuestionRepository.
func (p *postgresQuestionRepository) DeleteQuestion(questionID string) error {
	panic("unimplemented")
}

// FindQuestionByID implements QuestionRepository.
func (p *postgresQuestionRepository) FindQuestionByID(questionID string) (*models.Question, error) {
	panic("unimplemented")
}

// UpdateQuestion implements QuestionRepository.
func (p *postgresQuestionRepository) UpdateQuestion(question *models.Question) error {
	panic("unimplemented")
}

// FindQuestionByQuestionIdFeedbackId implements QuestionRepository.
func (p *postgresQuestionRepository) FindQuestionByQuestionIdFeedbackId(questionID uint, feedbackID string) (models.Question, error) {
	var matchingQuestions models.Question

	res := db.Where("id = ? AND feedback_id = ?", questionID, feedbackID).First(&matchingQuestions)
	return matchingQuestions, res.Error
}

// GetQuestionsByFeedbackID implements QuestionRepository.
func (p *postgresQuestionRepository) GetQuestionsByFeedbackID(feedbackID string) ([]models.Question, error) {
	var questions []models.Question
	var questions1 models.Question

	res := db.Find(&questions1)
	return questions, res.Error
}

// InsertQuestion implements QuestionRepository.
func (p *postgresQuestionRepository) InsertQuestion(question *models.Question) error {
	res := db.Create(&question)
	return res.Error
}

func newPostgresQuestionRepository(db *gorm.DB) QuestionRepository {
	return &postgresQuestionRepository{
		postgresdb: db,
	}
}
func GetQuestionRepository() QuestionRepository {
	return newPostgresQuestionRepository(db)
}
