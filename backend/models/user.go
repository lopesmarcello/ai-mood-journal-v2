package models

import "time"

type User struct {
	ID               uint      `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	PasswordHash     string    `json:"-"`
	IsPro            bool      `json:"is_pro"`
	TrialEntriesUsed int       `json:"trials_entries_used"`
	CreatedAt        time.Time `json:"created_at"`
}
