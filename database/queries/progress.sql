
-- name: CreateProgress :one
INSERT INTO progresses ( status) VALUES ( @status ) RETURNING progress_id;

-- name: GetProgress :one
SELECT status FROM progresses WHERE progress_id = @progress_id;

-- TODO: add a way to do cleanup of old progresses
