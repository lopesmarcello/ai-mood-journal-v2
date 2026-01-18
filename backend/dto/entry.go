package dto

type CreateEntryRequest struct {
	Content string `json:"content" binding:"required"`
}
