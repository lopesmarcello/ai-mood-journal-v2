package handlers

import (
	"net/http"
	"strconv"

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

	insightResponse := &dto.AIInsightResponse{
		Summary:    insight.Summary,
		Themes:     insight.Themes,
		Feelings:   insight.Feelings,
		Reflection: insight.FollowUp,
	}

	entryResponse := &dto.CreateEntryResponse{
		Content:   entry.Content,
		CreatedAt: entry.CreatedAt.Time,
		ID:        entry.ID,
		UserID:    entry.UserID.Int32,
	}

	c.JSON(http.StatusCreated, gin.H{
		"entry":   entryResponse,
		"insight": insightResponse,
	})
}

func (h *JournalHandler) List(c *gin.Context) {
	userID := int32(c.GetUint("user_id"))

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	entries, hasMore, err := h.service.ListEntries(c.Request.Context(), userID, int32(page))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch entries"})
		return
	}

	entriesResponse := make([]dto.CreateEntryResponse, 0, len(entries))
	for _, e := range entries {
		entriesResponse = append(entriesResponse,
			dto.CreateEntryResponse{
				Content:   e.Content,
				CreatedAt: e.CreatedAt.Time,
				ID:        e.ID,
				UserID:    e.UserID.Int32,
			},
		)
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data: entriesResponse,
		Pagination: dto.Pagination{
			CurrentPage: int32(page),
			PageSize:    10,
			HasMore:     hasMore,
		},
	})
}

func (h *JournalHandler) GetByID(c *gin.Context) {
	userID := int32(c.GetUint("user_id"))

	idStr := c.Param("id")
	entryID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry id"})
		return
	}

	entry, insight, err := h.service.GetEntryDetail(c.Request.Context(), userID, int32(entryID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"entry":   entry,
		"insight": insight,
	})
}
