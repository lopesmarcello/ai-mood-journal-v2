package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lopesmarcello/ai-journal/ai"
	"github.com/lopesmarcello/ai-journal/config"
	db "github.com/lopesmarcello/ai-journal/db/sqlc"
	"github.com/lopesmarcello/ai-journal/handlers"
	"github.com/lopesmarcello/ai-journal/middleware"
	"github.com/lopesmarcello/ai-journal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	defer pool.Close()

	queries := db.New(pool)

	authService := services.NewAuthService(queries, cfg.JWTSecret)

	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()

	prompt, err := ai.LoadSystemPrompt("prompts/therapist_v1.txt")
	if err != nil {
		log.Fatalf("Failed to load prompt: %w", err)
	}

	aiClient := ai.NewAIClient(cfg.RouteLLMKey, prompt)
	journalService := services.NewJournalService(pool, aiClient)
	journalHandler := handlers.NewJournalHandler(journalService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.POST("/auth/logout", authHandler.Me)

		protected.POST("/entries", journalHandler.Create)
	}

	log.Printf("Server starting on port %s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
