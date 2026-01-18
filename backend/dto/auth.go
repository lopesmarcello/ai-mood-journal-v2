package dto

import "github.com/golang-jwt/jwt/v5"

type AuthResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	IsPro bool   `json:"is_pro"`
}

type JWTClaims struct {
	UserID uint `json:"user_id"`
	IsPro  bool `json:"is_pro"`
	jwt.RegisteredClaims
}
