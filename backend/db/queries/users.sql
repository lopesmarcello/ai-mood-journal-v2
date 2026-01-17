-- name: CreateUser :one
INSERT INTO users (
  email, name, password_hash
) VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;


-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateTrialEntriesUsed :one
UPDATE users 
set trial_entries_used = trial_entries_used + 1
WHERE id = $1
RETURNING trial_entries_used;

-- name: SetUserPro :one
UPDATE users 
set is_pro =$2
WHERE id = $1
RETURNING *;
