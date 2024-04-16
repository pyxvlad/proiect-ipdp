
-- name: CreateSessionToken :exec
INSERT INTO sessions (account_id, token) VALUES (@account_id, @token);

-- name: GetSessionAccount :one
SELECT account_id FROM sessions WHERE token = @token LIMIT 1;

