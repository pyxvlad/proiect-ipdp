
-- name: CreateBook :one
INSERT INTO books (account_id, title, author, progress_id, publisher_id) VALUES ( @account_id, @title, @author, @progress_id, @publisher_id ) RETURNING book_id;

-- name: GetBooksWithStatuses :many
SELECT book_id, title, author, status FROM books JOIN progresses ON books.progress_id = progresses.progress_id WHERE account_id = @account_id;

-- name: GetBooksWithCoversAndStatuses :many
SELECT book_id, title, author, status, cover_hash FROM books JOIN progresses ON books.progress_id = progresses.progress_id WHERE account_id = @account_id;

-- name: LinkBookToProgress :exec
UPDATE books SET progress_id = @progress_id WHERE book_id = @book_id;

-- name: LinkBookToDuplicate :exec
UPDATE books SET duplicate_id = @duplicate_id WHERE book_id = @book_id;

-- name: GetDuplicateOfBook :one
SELECT duplicate_id FROM books WHERE book_id = @book_id;

-- name: SetBookCoverHash :exec
UPDATE books SET cover_hash = @cover_hash WHERE book_id = @book_id;

-- name: ChangePublisher :exec
UPDATE books SET publisher_id = @publisher_id WHERE book_id = @book_id;

