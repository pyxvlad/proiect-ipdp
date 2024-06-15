-- name: CreateSeries :one
INSERT INTO series ( name, account_id) VALUES ( @name, @account_id ) RETURNING series_id;

-- name: CreateCollection :one
INSERT INTO collections ( name, account_id) VALUES ( @name, @account_id ) RETURNING collection_id;

-- name: AddBookToSeries :exec
INSERT INTO book_series (series_id, book_id, volume) VALUES (@series_id, @book_id, @volume);

-- name: AddBookToCollection :exec
INSERT INTO book_collections (collection_id, book_id, book_number) VALUES (@collection_id, @book_id, @book_number);

-- name: RemoveBookFromSeries :exec
DELETE FROM book_series WHERE book_id = @book_id AND series_id = @series_id;

-- name: RemoveBookFromCollection :exec
DELETE FROM book_collections WHERE book_id = @book_id AND collection_id = @collection_id;

-- name: DeleteSeries :exec
DELETE FROM series WHERE series_id = @series_id AND account_id = @account_id;

-- name: DeleteCollection :exec
DELETE FROM collections WHERE collection_id = @collection_id AND account_id = @account_id;

-- name: ListSeriesForAccount :many
SELECT series_id, name FROM series WHERE account_id = @accoount_id;

-- name: ListCollectionsForAccount :many
SELECT collection_id, name FROM collections WHERE account_id = @accoount_id;

-- name: RenameSeries :exec
UPDATE series SET name = @name
WHERE series_id = @series_id AND account_id = @account_id;

-- name: RenameCollection :exec
UPDATE collections SET name = @name
WHERE collection_id = @collection_id AND account_id = @account_id;

-- name: GetNameOfSeries :one
SELECT name
FROM series
WHERE series_id = @series_id AND account_id = @account_id;

-- name: GetNameOfCollection :one
SELECT name
FROM collections
WHERE collection_id = @collection_id AND account_id = @account_id;

-- name: GetSeriesForBook :one
SELECT name, volume, series.series_id
FROM book_series JOIN series
ON book_series.series_id = series.series_id
WHERE account_id = @account_id AND book_id = @book_id;

-- name: GetCollectionForBook :one
SELECT name, book_number, collections.collection_id
FROM book_collections JOIN collections
ON book_collections.collection_id = collections.collection_id
WHERE account_id = @account_id AND book_id = @book_id;

-- name: ChangeBookVolumeInSeries :exec
UPDATE book_series
SET volume = @volume
WHERE book_id = @book_id;

-- name: ChangeBookNumberInCollection :exec
UPDATE book_collections
SET book_number = @book_number
WHERE book_id = @book_id;
