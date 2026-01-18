package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	db "github.com/lopesmarcello/ai-journal/db/sqlc"
	"github.com/lopesmarcello/ai-journal/dto"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries   *db.Queries
	jwtSecret string
}

func NewAuthService(queries *db.Queries, jwtSecret string) *AuthService {
	return &AuthService{
		queries:   queries,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(uint(user.ID), user.IsPro.Bool)
	if err != nil {
		return nil, "", err
	}

	return &dto.AuthResponse{
		ID:    uint(user.ID),
		Email: user.Email,
		Name:  user.Name,
		IsPro: user.IsPro.Bool,
	}, token, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, string, error) {
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.GenerateToken(uint(user.ID), user.IsPro.Bool)
	if err != nil {
		return nil, "", err
	}

	return &dto.AuthResponse{
		ID:    uint(user.ID),
		Email: user.Email,
		Name:  user.Name,
		IsPro: user.IsPro.Bool,
	}, token, nil
}

func (s *AuthService) GenerateToken(userID uint, isPro bool) (string, error) {
	dayInHours := 24 * time.Hour
	aMonthFromNow := time.Now().Add(30 * dayInHours)

	claims := dto.JWTClaims{
		UserID: userID,
		IsPro:  isPro,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(aMonthFromNow),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
