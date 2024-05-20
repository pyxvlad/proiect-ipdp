
-- name: CreateDuplicate :one
INSERT INTO duplicates (always_null) VALUES (NULL) RETURNING duplicate_id;

