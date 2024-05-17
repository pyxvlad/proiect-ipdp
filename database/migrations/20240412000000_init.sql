
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
CREATE TABLE publishers (
	publisher_id INTEGER PRIMARY KEY,
	account_id INTEGER NOT NULL,

	name text NOT NULL
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE books (
	book_id INTEGER PRIMARY KEY,
	account_id INTEGER NOT NULL,

	title text NOT NULL,
	author text NOT NULL,

	cover_hash text,

	duplicate_id INTEGER DEFAULT NULL,
	progress_id INTEGER NOT NULL,
	publisher_id INTEGER NOT NULL,

	FOREIGN KEY(account_id) REFERENCES accounts(account_id),
	FOREIGN KEY(duplicate_id) REFERENCES duplicates(duplicate_id),
	FOREIGN KEY(progress_id) REFERENCES progresses(progress_id),
	FOREIGN KEY(publisher_id) REFERENCES publishers(publisher_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE series (
	series_id INTEGER PRIMARY KEY,
	account_id INTEGER,

	name text NOT NULL,

	FOREIGN KEY(account_id) REFERENCES accounts(account_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE book_series (
	series_id INTEGER,
	book_id INTEGER,

	volume INTEGER NOT NULL,

	PRIMARY KEY(series_id, book_id),
	FOREIGN KEY(book_id) REFERENCES books(book_id),
	FOREIGN KEY(series_id) REFERENCES series(series_id),

	CHECK ( volume is null OR volume > 0 )
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE collections (
	collection_id INTEGER PRIMARY KEY,
	account_id INTEGER,

	name text NOT NULL,

	FOREIGN KEY(account_id) REFERENCES accounts(account_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE book_collections (
	collection_id INTEGER,
	book_id INTEGER,

	book_number INTEGER,

	PRIMARY KEY(collection_id, book_id),
	FOREIGN KEY(book_id) REFERENCES books(book_id),
	FOREIGN KEY(collection_id) REFERENCES collections(collection_id),


	CHECK ( book_number > 0 )
);
-- +goose StatementEnd

