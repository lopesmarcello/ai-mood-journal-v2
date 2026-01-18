package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lopesmarcello/ai-journal/dto"
	"github.com/lopesmarcello/ai-journal/services"
)

type JournalHandler struct {
	service *services.JournalService
}

func NewJournalHandler(s *services.JournalService) *JournalHandler {
	return &JournalHandler{service: s}
}

func (h *JournalHandler) Create(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context missing"})
	}
	userID := int32(val.(uint))

	var req dto.CreateEntryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	entry, insight, err := h.service.CreateEntry(c.Request.Context(), userID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save entry"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"entry":   entry,
		"insight": insight,
	})
}
