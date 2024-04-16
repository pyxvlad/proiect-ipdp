
-- name: CreateAccountWithEmail :one
INSERT INTO accounts (email, password) VALUES (@email, @password) RETURNING account_id;

-- name: GetPasswordByEmail :one
SELECT account_id, password FROM accounts WHERE email = @email;

-- name: GetAccountByEmail :one
SELECT account_id FROM accounts WHERE email = @email;



-- name: GetPasswordForAccount :one
SELECT password FROM accounts WHERE account_id = @account_id;

-- name: SetPasswordForAccount :exec
UPDATE accounts SET password = @password WHERE account_id = @account_id;

-- name: FindAccountByID :one
SELECT * FROM accounts WHERE account_id = @account_id;

-- name: FindAccountByEmail :one
SELECT * FROM accounts WHERE email = @email;
