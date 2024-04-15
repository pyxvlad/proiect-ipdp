
-- +goose Up

-- +goose StatementBegin
CREATE TABLE accounts (
	account_id INTEGER PRIMARY KEY,
	email text NOT NULL,
	password text NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE sessions (
	session_id INTEGER PRIMARY KEY,
	account_id INTEGER NOT NULL,
	token text NOT NULL,

	FOREIGN KEY(account_id) REFERENCES accounts(account_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE progresses (
	progress_id INTEGER PRIMARY KEY,
	status text CHECK ( status IN ('to be read', 'in progress', 'read', 'dropped', 'uncertain') ) NOT NULL DEFAULT 'to be read'
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE duplicates (
	duplicate_id INTEGER PRIMARY KEY,

	always_null INTEGER DEFAULT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE books (
	book_id INTEGER PRIMARY KEY,
	account_id INTEGER NOT NULL,

	title text NOT NULL,
	author text NOT NULL,

	duplicate_id INTEGER DEFAULT NULL,
	progress_id INTEGER NOT NULL,

	FOREIGN KEY(account_id) REFERENCES accounts(account_id),
	FOREIGN KEY(duplicate_id) REFERENCES duplicates(duplicate_id),
	FOREIGN KEY(progress_id) REFERENCES progress(progress_id)
);
-- +goose StatementEnd


