package models

// Recommendation represents a heuristic recommendation for the best model
// across a multi-model evaluation run.
type Recommendation struct {
	RecommendedModel string                `json:"recommended_model"`
	HeuristicScore   float64               `json:"heuristic_score"`
	Reason           string                `json:"reason"`
	WinnerMarginPct  float64               `json:"winner_margin_pct"`
	Weights          RecommendationWeights `json:"weights"`
	ModelScores      []ModelScore          `json:"all_models"`
}

// RecommendationWeights defines the weighting scheme for heuristic scoring.
type RecommendationWeights struct {
	AggregateScore float64 `json:"aggregate_score"`
	PassRate       float64 `json:"pass_rate"`
	Consistency    float64 `json:"consistency"`
	Speed          float64 `json:"speed"`
}

// ModelScore holds the heuristic score and rank for a single model.
type ModelScore struct {
	ModelID        string             `json:"model_id"`
	HeuristicScore float64            `json:"heuristic_score"`
	Rank           int                `json:"rank"`
	Scores         map[string]float64 `json:"component_scores,omitempty"`
}
