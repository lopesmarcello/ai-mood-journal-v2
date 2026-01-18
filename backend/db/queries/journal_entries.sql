-- name: CreateJournalEntry :one
INSERT INTO journal_entries (
  user_id, content
) VALUES ($1, $2)
RETURNING *;

-- name: ListEntriesByUser :many
SELECT * FROM journal_entries
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetSingleEntryByIDs :one
SELECT * FROM journal_entries
WHERE user_id = $1 AND id = $2 LIMIT 1;

-- name: CountEntrieByUser :one
SELECT COUNT(*) FROM journal_entries
WHERE user_id = $1;

