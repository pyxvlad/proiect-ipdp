-- name: CreatePublisher :one
INSERT INTO publishers ( account_id, name) VALUES ( @account_id, @name)
RETURNING publisher_id;

-- name: ListPublishersForAccount :many
SELECT publisher_id, name FROM publishers WHERE account_id = @account_id;

-- name: DeletePublisher :exec
DELETE FROM publishers
WHERE account_id = @account_id AND publisher_id = @publisher_id;

-- name: RenamePublisher :exec
UPDATE publishers SET name = @name
WHERE publisher_id = @publisher_id AND account_id = @account_id;
