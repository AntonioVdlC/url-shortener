CREATE TABLE IF NOT EXISTS links (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hash VARCHAR(8),
  url TEXT,
  is_safe BOOLEAN,
  is_safe_next_check_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS hash_idx on links (hash);
