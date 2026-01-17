package models

import "time"

type JournalEntry struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	Insight   AIInsight `json:"insight"`
	CreatedAt time.Time `json:"created_at"`
}
