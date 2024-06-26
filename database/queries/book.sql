
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

-- name: GetAllBookData :one
SELECT book_id, title, author, status, cover_hash, publisher_id, duplicate_id
FROM books JOIN progresses ON books.progress_id = progresses.progress_id
WHERE account_id = @account_id AND book_id = @book_id;

-- name: GetDuplicatesForBook :many
SELECT book_id, title, author, status, cover_hash
FROM books JOIN progresses ON books.progress_id = progresses.progress_id
WHERE account_id = @account_id AND duplicate_id = @duplicate_id AND book_id != @book_id;

-- name: SetBookTitle :exec
UPDATE books SET title = @title
WHERE book_id = @book_id AND @account_id = @account_id;

-- name: SetBookAuthor :exec
UPDATE books SET author = @author
WHERE book_id = @book_id AND @account_id = @account_id;


-- name: GetBooksForPublisher :many
SELECT book_id, title, author, status, cover_hash
FROM books JOIN progresses ON books.progress_id = progresses.progress_id
WHERE publisher_id = @publisher_id;

-- name: GetBooksForAuthor :many
SELECT book_id, title, author, status, cover_hash
FROM books JOIN progresses ON books.progress_id = progresses.progress_id
WHERE author = @author;

-- name: GetBooksForCollection :many
SELECT books.book_id, title, author, status, cover_hash, book_number
FROM books JOIN progresses ON books.progress_id = progresses.progress_id
JOIN book_collections ON books.book_id = book_collections.book_id
WHERE collection_id = @collection_id;

-- name: GetBooksForSeries :many
SELECT books.book_id, title, author, status, cover_hash, volume
FROM books JOIN progresses ON books.progress_id = progresses.progress_id
JOIN book_series ON books.book_id = book_series.book_id
WHERE series_id = @series_id;



-- name: GetBooksForSlice :many
SELECT book_id, title, author, status, cover_hash
FROM books JOIN progresses ON books.progress_id = progresses.progress_id
WHERE book_id IN (sqlc.slice('book_ids'));


