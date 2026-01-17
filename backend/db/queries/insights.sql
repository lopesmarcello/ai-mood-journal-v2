-- name: CreateInsight :one
INSERT INTO insights (
  entry_id, themes, feelings, summary, follow_up
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetInsightByEntryID :one
SELECT * from insights 
WHERE entry_id = $1 LIMIT 1;
