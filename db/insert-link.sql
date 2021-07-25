INSERT INTO links (
  hash,
  url,
  is_safe,
  is_safe_next_check_at
) VALUES ($1, $2, $3, $4);
