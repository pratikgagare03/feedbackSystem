package models

type Stats struct {
	Question     string       `json:"question"`
	QuestionType QuestionType `json:"question_type"`
	Stats        interface{}  `json:"stats"`
}
type RatingsStats struct {
	MaxRatingsRange   int                `json:"max_ratings_range"`
	TotalResponses    int                `json:"total_responses"`
	RatingsCount      map[string]int     `json:"ratings_count"`
	RatingsPercentage map[string]float64 `json:"ratings_percentage"`
	AverageRating     float64            `json:"average_rating"`
}

type MCQStats struct {
	AllOptions        []string           `json:"all_options"`
	TotalResponses    int                `json:"total_responses"`
	OptionsCount      map[string]int     `json:"options_count"`
	OptionsPercentage map[string]float64 `json:"options_percentage"`
}

type TextInputStats struct {
	TotalResponses int      `json:"total_responses"`
	Answers        []string `json:"answers"`
}
