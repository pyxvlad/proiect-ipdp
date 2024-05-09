
-- name: CreateBook :one
INSERT INTO books (account_id, title, author, progress_id) VALUES ( @account_id, @title, @author, @progress_id ) RETURNING book_id;

-- name: GetBooksWithStatuses :one
SELECT title, author, status FROM books JOIN progresses ON books.progress_id = progresses.progress_id WHERE account_id = @account_id;

-- name: LinkBookToProgress :exec
UPDATE books SET progress_id = @progress_id WHERE book_id = @book_id;

-- name: LinkBookToDuplicate :exec
UPDATE books SET duplicate_id = @duplicate_id WHERE book_id = @book_id;

-- name: GetDuplicateOfBook :one
SELECT duplicate_id FROM books WHERE book_id = @book_id;

