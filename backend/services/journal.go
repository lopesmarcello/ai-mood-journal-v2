package services

import (
	"context"
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
	if err != nil {
		return nil, nil, err
	}

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
			slog.Error("AI Insight failed:", "error", aiErr)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return &entry, insight, nil
}

func (s *JournalService) ListEntries(ctx context.Context, userID int32, page int32) ([]db.JournalEntry, bool, error) {
	limit := int32(10)
	offset := (page - 1) * limit

	entries, err := s.queries.ListEntriesByUser(ctx, db.ListEntriesByUserParams{
		UserID: pgtype.Int4{Int32: userID, Valid: true},
		Limit:  limit + 1,
		Offset: offset,
	})
	if err != nil {
		return nil, false, err
	}

	hasMore := false
	if len(entries) > int(limit) {
		hasMore = true
		entries = entries[:limit]
	}

	return entries, hasMore, nil
}

func (s *JournalService) GetEntryDetail(ctx context.Context, userID int32, entryID int32) (*db.JournalEntry, *db.Insight, error) {
	entry, err := s.queries.GetSingleEntryByIDs(ctx, db.GetSingleEntryByIDsParams{
		UserID: pgtype.Int4{Int32: userID, Valid: true},
		ID:     entryID,
	})
	if err != nil {
		return nil, nil, err
	}

	insight, err := s.queries.GetInsightByEntryID(ctx, pgtype.Int4{
		Int32: entryID,
		Valid: true,
	})
	if err != nil {
		return &entry, nil, err
	}

	return &entry, &insight, nil
}
