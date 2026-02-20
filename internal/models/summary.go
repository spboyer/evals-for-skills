package models

import "time"

// MultiSkillSummary aggregates results across multiple skill evaluations.
type MultiSkillSummary struct {
	Timestamp time.Time      `json:"timestamp"`
	Skills    []SkillSummary `json:"skills"`
	Overall   OverallSummary `json:"overall"`
}

// SkillSummary contains aggregated metrics for a single skill evaluation.
type SkillSummary struct {
	SkillName      string   `json:"skill_name"`
	Models         []string `json:"models"`
	PassRate       float64  `json:"pass_rate"`
	AggregateScore float64  `json:"aggregate_score"`
	OutputFiles    []string `json:"output_files"`
}

// OverallSummary contains cross-skill aggregated metrics.
type OverallSummary struct {
	TotalSkills       int     `json:"total_skills"`
	TotalModels       int     `json:"total_models"`
	AvgPassRate       float64 `json:"avg_pass_rate"`
	AvgAggregateScore float64 `json:"avg_aggregate_score"`
}
