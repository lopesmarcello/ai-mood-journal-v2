package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lopesmarcello/ai-journal/ai"
	db "github.com/lopesmarcello/ai-journal/db/sqlc"
)

type JournalService struct {
	pool     *pgxpool.Pool
	queries  *db.Queries
	aiClient *ai.AIClient
}

func NewJournalService(pool *pgxpool.Pool, aiClient *ai.AIClient) *JournalService {
	return &JournalService{
		pool:     pool,
		queries:  db.New(pool),
		aiClient: aiClient,
	}
}

func (s *JournalService) CreateEntry(ctx context.Context, userID int32, content string) (*db.JournalEntry, *db.Insight, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	user, err := qtx.GetUserByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	entry, err := qtx.CreateJournalEntry(ctx, db.CreateJournalEntryParams{
		UserID:  pgtype.Int4{Int32: userID, Valid: true},
		Content: content,
	})

	var insight *db.Insight

	maxTrialEntries := 10
	userConditionToHaveInsight := user.IsPro.Bool || (user.TrialEntriesUsed.Int32 < int32(maxTrialEntries))

	if userConditionToHaveInsight {
		aiResp, aiErr := s.aiClient.GenerateInsight(ctx, content)
		if aiErr == nil {
			savedInsight, err := qtx.CreateInsight(ctx, db.CreateInsightParams{
				EntryID:  pgtype.Int4{Int32: entry.ID, Valid: true},
				Summary:  aiResp.Summary,
				Themes:   aiResp.Themes,
				Feelings: aiResp.Feelings,
				FollowUp: aiResp.Reflection,
			})

			if err == nil {
				insight = &savedInsight
				if !user.IsPro.Bool {
					_, err = qtx.UpdateTrialEntriesUsed(ctx, userID)
					if err != nil {
						slog.Error("Error updating TrialEntriesUsed", "error", err)
					}
				}
			}
		} else {
			fmt.Printf("AI Insight failed: %w", aiErr)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return &entry, insight, nil
}
