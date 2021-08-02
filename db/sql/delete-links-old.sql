DELETE 
FROM links 
WHERE created_at < current_timestamp - interval '30 days';
