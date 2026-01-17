-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS journal_entries (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

---- create above / drop below ----
DROP TABLE IF EXISTS journal_entries;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
