package models

type RatingsStats struct {
	Question          string          `json:"question"`
	QuestionType      string          `json:"question_type"`
	TotalResponses    int             `json:"total_responses"`
	RatingsCount      map[int]int     `json:"ratings_count"`
	RatingsPercentage map[int]float64 `json:"ratings_percentage"`
	AverageRating     float64         `json:"average_rating"`
}

type MCQStats struct {
	Question          string             `json:"question"`
	QuestionType      string             `json:"question_type"`
	OptionsCount      map[string]int     `json:"options_count"`
	OptionsPercentage map[string]float64 `json:"options_percentage"`
}

type TextInputStats struct {
	Question     string   `json:"question"`
	QuestionType string   `json:"question_type"`
	Answers      []string `json:"answers"`
}
