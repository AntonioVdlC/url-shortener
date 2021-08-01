CREATE TABLE IF NOT EXISTS links (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  hash VARCHAR(8),
  url TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS hash_idx on links (hash);
CREATE UNIQUE INDEX IF NOT EXISTS created_at_idx on links (created_at);
