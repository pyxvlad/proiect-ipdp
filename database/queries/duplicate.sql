
-- name: CreateDuplicate :exec
INSERT INTO duplicates (always_null) VALUES (NULL) RETURNING duplicate_id;


