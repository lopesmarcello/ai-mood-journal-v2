package models

type AIInsight struct {
	Themes            []string `json:"themes"`
	PrimaryFeelings   []string `json:"feelings"`
	ReflectiveSummary string   `json:"summary"`
	FollowUpQuestion  string   `json:"question"`
}
