package dto

import "time"

type CreateEntryRequest struct {
	Content string `json:"content" binding:"required"`
}

type CreateEntryResponse struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
}
