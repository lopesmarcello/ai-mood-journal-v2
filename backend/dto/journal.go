package dto

type AIInsightResponse struct {
	Summary    string   `json:"summary"`
	Themes     []string `json:"themes"`
	Feelings   []string `json:"feelings"`
	Reflection string   `json:"reflection"`
}
