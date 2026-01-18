package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/lopesmarcello/ai-journal/db/sqlc"
	"github.com/lopesmarcello/ai-journal/dto"
	"github.com/lopesmarcello/ai-journal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: s,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	h.setAuthCookie(c, token)
	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}

	h.setAuthCookie(c, token)
	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) TogglePro(c *gin.Context) {
	userID := c.GetUint("user_id")
	isPro := c.GetBool("is_pro")

	var newUser db.User

	if isPro {
		foundUser, err := h.authService.DowngradePro(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		newUser = foundUser
	} else {
		foundUser, err := h.authService.UpgradeToPro(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		newUser = foundUser
	}

	c.JSON(http.StatusOK, gin.H{"user_id": newUser.ID, "is_pro": newUser.IsPro})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	isPro, _ := c.Get("is_pro")

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"is_pro":  isPro,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
	// name, value, maxAge (seconds), path, domain, secure, httpOnly
	c.SetCookie("auth_token", token, 3600*24*30, "/", "", false, true)
}
