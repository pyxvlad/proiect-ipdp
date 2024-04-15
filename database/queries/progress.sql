
-- name: CreateProgress :one
INSERT INTO progresses ( status) VALUES ( @status ) RETURNING progress_id;

-- name: GetProgress :one
SELECT status FROM progresses WHERE progress_id = @progress_id;

-- name: SetBookStatus :exec
UPDATE progresses SET status = @status WHERE progress_id in (
	SELECT progress_id FROM books WHERE book_id = @book_id
);
-- TODO: make SetBookStatus a CTE using "WITH"

-- TODO: add a way to do cleanup of old progresses
