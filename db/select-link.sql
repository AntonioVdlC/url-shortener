SELECT url 
FROM links 
WHERE is_safe = TRUE 
  AND hash = $1;
