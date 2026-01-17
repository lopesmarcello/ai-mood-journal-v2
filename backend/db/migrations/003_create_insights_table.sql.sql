-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS insights (
  id SERIAL PRIMARY KEY,
  entry_id INT UNIQUE REFERENCES journal_entries(id) ON DELETE CASCADE,
  themes TEXT[],
  feelings TEXT[],
  summary TEXT NOT NULL,
  follow_up TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

---- create above / drop below ----
DROP TABLE IF NOT EXISTS insights;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
