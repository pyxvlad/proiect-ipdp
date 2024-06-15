package services

import (
	"cmp"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path"
	"slices"

	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/rs/zerolog"
)

type BookService struct {
	ImagePath string
}

func NewBookService(imagePath string) BookService {
	return BookService{
		ImagePath: imagePath,
	}
}

// CreateBook creates a new book.
func (b *BookService) CreateBook(
	ctx context.Context,
	account types.AccountID,
	title string,
	author string,
	status types.Status,
	publisherID types.PublisherID,
) (types.BookID, error) {
	log := zerolog.Ctx(ctx)

	tx, err := DBTxLike(ctx).Begin()
	defer tx.Rollback()
	if err != nil {
		log.Err(err).Msg("while trying to begin transaction")
		return types.InvalidBookID, err
	}

	queries := database.New(tx)

	progressID, err := queries.CreateProgress(ctx, status)
	if err != nil {
		log.Err(err).Msg("while trying to create progress for book")
		return types.InvalidBookID, err
	}

	bookID, err := queries.CreateBook(ctx, database.CreateBookParams{
		AccountID:   account,
		Title:       title,
		Author:      author,
		ProgressID:  progressID,
		PublisherID: publisherID,
	})

	if err != nil {
		log.Err(err).Msg("while trying to create book")
		return types.InvalidBookID, err
	}

	err = tx.Commit()
	if err != nil {
		log.Err(err).Msg("while trying to commit transaction")
		return types.InvalidBookID, err
	}

	return bookID, nil
}

// SetBookStatus sets the status of a book.
func (b *BookService) SetBookStatus(
	ctx context.Context, bookID types.BookID, status types.Status,
) error {
	err := DB(ctx).SetBookStatus(ctx, database.SetBookStatusParams{
		Status: status,
		BookID: bookID,
	})
	if err != nil {
		return err
	}
	return nil
}

// MarkBookAsDuplicate marks a duplicatedBookID as a duplicate of bookID.
func (b *BookService) MarkBookAsDuplicate(ctx context.Context, bookID types.BookID, duplicatedBookID types.BookID) error {
	log := zerolog.Ctx(ctx).With().Caller().Int("bookID", int(bookID)).Int("duplicatedBookID", int(duplicatedBookID)).Logger()
	txlike := DBTxLike(ctx)

	tx, err := txlike.Begin()
	if err != nil {
		log.Err(err).Msg("while trying to begin the transaction")
		return err
	}

	defer tx.Rollback()

	queries := database.New(tx)
	optDuplicateID, err := queries.GetDuplicateOfBook(ctx, bookID)
	if err != nil {

		log.Err(err).Msg("while trying to query for a duplicateID of bookID")
		return err
	}

	if optDuplicateID.Valid {
		err = queries.LinkBookToDuplicate(ctx, database.LinkBookToDuplicateParams{
			DuplicateID: optDuplicateID,
			BookID:      bookID,
		})
		if err != nil {
			log.Err(err).Msg("while trying to link book to duplicate")
			return err
		}
	} else {
		var duplicateID int64
		duplicateID, err = queries.CreateDuplicate(ctx)
		if err != nil {
			log.Err(err).Msg("while trying to create a new duplicate")
			return err
		}

		err = queries.LinkBookToDuplicate(ctx, database.LinkBookToDuplicateParams{
			DuplicateID: sql.NullInt64{Int64: duplicateID, Valid: true},
			BookID:      bookID,
		})
		if err != nil {
			log.Err(err).Msg("while trying to link book bookID to duplicate")
			return err
		}

		err = queries.LinkBookToDuplicate(ctx, database.LinkBookToDuplicateParams{
			DuplicateID: sql.NullInt64{Int64: duplicateID, Valid: true},
			BookID:      duplicatedBookID,
		})
		if err != nil {
			log.Err(err).Msg("while trying to link book duplicated to duplicate")
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}
	return nil
}

type BookData struct {
	BookID types.BookID
	Title  string
	Author string
	Status types.Status
}

func (b *BookService) ListBooksForAccount(ctx context.Context, accountID types.AccountID) ([]BookData, error) {
	rows, err := DB(ctx).GetBooksWithStatuses(ctx, accountID)
	if err != nil {
		return nil, err
	}

	data := make([]BookData, 0, len(rows))
	for _, row := range rows {
		var bd BookData
		bd.Title = row.Title
		bd.Author = row.Author
		bd.Status = row.Status
		bd.BookID = row.BookID
		data = append(data, bd)
	}

	return data, nil
}

func (b *BookService) SetBookCover(ctx context.Context, bookID types.BookID, cover_image io.Reader) error {
	hasher := sha256.New()
	tee := io.TeeReader(cover_image, hasher)
	file, err := os.CreateTemp("/tmp/ipdp-img", "")
	if err != nil {
		return err
	}
	_, err = io.Copy(file, tee)
	if err != nil {
		return err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	err = os.Rename(file.Name(), path.Join(b.ImagePath, hash))
	if err != nil {
		return err
	}

	err = DB(ctx).SetBookCoverHash(ctx, database.SetBookCoverHashParams{
		CoverHash: sql.NullString{
			String: hash,
			Valid:  true,
		},
		BookID: bookID,
	})

	return err
}

type BookDataWithCovers struct {
	BookID    types.BookID
	Title     string
	Author    string
	Status    types.Status
	CoverHash string
}

func (b *BookService) ListBooksWithCoversForAccount(
	ctx context.Context, accountID types.AccountID,
) ([]BookDataWithCovers, error) {
	rows, err := DB(ctx).GetBooksWithCoversAndStatuses(ctx, accountID)
	if err != nil {
		return nil, err
	}

	data := make([]BookDataWithCovers, 0, len(rows))
	for _, row := range rows {
		var bd BookDataWithCovers
		bd.Title = row.Title
		bd.Author = row.Author
		bd.Status = row.Status
		bd.BookID = row.BookID
		if row.CoverHash.Valid {
			bd.CoverHash = row.CoverHash.String
		} else {
			bd.CoverHash = ""
		}
		data = append(data, bd)
	}

	return data, nil

}

func (b *BookService) SetBookPublisher(
	ctx context.Context, bookID types.BookID, publisherID types.PublisherID,
) error {
	return DB(ctx).ChangePublisher(ctx, database.ChangePublisherParams{
		PublisherID: publisherID,
		BookID:      bookID,
	})
}

type BookAllData struct {
	BookID    types.BookID
	Title     string
	Author    string
	Status    types.Status
	CoverHash string

	PublisherName string
	Collection    struct {
		Name   string
		Number uint
	}
	Series struct {
		Name   string
		Volume uint
	}
	Duplicates []BookDataWithCovers
}

func (b *BookService) GetAllDataForBook(
	ctx context.Context, bookID types.BookID, accountID types.AccountID,
) (BookAllData, error) {
	var bookData BookAllData
	log := zerolog.Ctx(ctx).With().Caller().Logger()
	db := DB(ctx)
	rowBookData, err := db.GetAllBookData(ctx, database.GetAllBookDataParams{
		AccountID: accountID,
		BookID:    bookID,
	})

	if err != nil {
		log.Err(err).Msg("while trying to get all book data")
		return BookAllData{}, err
	}

	bookData.BookID = rowBookData.BookID
	bookData.Author = rowBookData.Author
	bookData.Title = rowBookData.Title
	if rowBookData.CoverHash.Valid {
		bookData.CoverHash = rowBookData.CoverHash.String
	} else {
		bookData.CoverHash = ""
	}

	bookData.Status = rowBookData.Status

	collectionRow, err := db.GetCollectionForBook(ctx, database.GetCollectionForBookParams{
		AccountID: accountID,
		BookID:    bookID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Err(err).Msg("while trying to get collection for book")
		return BookAllData{}, err
	}

	bookData.Collection.Name = collectionRow.Name
	bookData.Collection.Number = uint(collectionRow.BookNumber.Int64)

	seriesRow, err := db.GetSeriesForBook(ctx, database.GetSeriesForBookParams{
		AccountID: accountID,
		BookID:    bookID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Err(err).Msg("while trying to get series for book")
		return BookAllData{}, err
	}

	bookData.Series.Name = seriesRow.Name
	bookData.Series.Volume = seriesRow.Volume

	bookData.PublisherName, err = db.GetNameOfPublisher(ctx, database.GetNameOfPublisherParams{
		AccountID:   accountID,
		PublisherID: rowBookData.PublisherID,
	})

	if err != nil {
		return BookAllData{}, err
	}

	if rowBookData.DuplicateID.Valid {
		duplicates, err := db.GetDuplicatesForBook(ctx, database.GetDuplicatesForBookParams{
			AccountID:   accountID,
			DuplicateID: rowBookData.DuplicateID,
			BookID:      bookID,
		})
		if err != nil {
			log.Err(err).Msg("while trying to get duplicates for book")
			return BookAllData{}, err
		}

		bookData.Duplicates = make([]BookDataWithCovers, 0, len(duplicates))
		for _, dup := range duplicates {
			var dupData BookDataWithCovers
			dupData.Title = dup.Title
			dupData.Author = dup.Author
			dupData.BookID = dup.BookID
			dupData.Status = dup.Status
			if dup.CoverHash.Valid {
				dupData.CoverHash = dup.CoverHash.String
			}

			bookData.Duplicates = append(bookData.Duplicates, dupData)
		}
	}

	return bookData, nil
}

func (bs *BookService) SetBookTitle(
	ctx context.Context,
	account_id types.AccountID,
	book_id types.BookID,
	title string,
) error {
	return DB(ctx).SetBookTitle(ctx, database.SetBookTitleParams{
		Title:     title,
		BookID:    book_id,
		AccountID: account_id,
	})
}

func (bs *BookService) SetBookAuthor(
	ctx context.Context,
	account_id types.AccountID,
	book_id types.BookID,
	author string,
) error {
	return DB(ctx).SetBookAuthor(ctx, database.SetBookAuthorParams{
		Author:    author,
		BookID:    book_id,
		AccountID: account_id,
	})
}

type SortFunc func(BookDataWithCovers, BookDataWithCovers) int

func SortByTitle(x, y BookDataWithCovers) int {
	return cmp.Compare(x.Title, y.Title)
}

func SortByAuthor(x, y BookDataWithCovers) int {
	return cmp.Compare(x.Author, y.Author)
}

func SortByStatus(x, y BookDataWithCovers) int {
	return cmp.Compare(x.Status, y.Status)
}

func (bs *BookService) ListBooksSorted(ctx context.Context, accountID types.AccountID, criteria SortFunc) ([]BookDataWithCovers, error) {
	books, err := bs.ListBooksWithCoversForAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(books, criteria)
	return books, nil
}

func (bs *BookService) FilterByAuthor(ctx context.Context, accountID types.AccountID, author string) ([]BookDataWithCovers, error) {
	rows, err := DB(ctx).GetBooksForAuthor(ctx, author)
	if err != nil {
		return nil, err
	}

	data := make([]BookDataWithCovers, 0, len(rows))
	for _, row := range rows {
		var bd BookDataWithCovers
		bd.Title = row.Title
		bd.Author = row.Author
		bd.Status = row.Status
		bd.BookID = row.BookID
		if row.CoverHash.Valid {
			bd.CoverHash = row.CoverHash.String
		} else {
			bd.CoverHash = ""
		}
		data = append(data, bd)
	}

	return data, nil
}

func (bs *BookService) FilterByPublisher(ctx context.Context, accountID types.AccountID, publisherID types.PublisherID) ([]BookDataWithCovers, error) {
	rows, err := DB(ctx).GetBooksForPublisher(ctx, publisherID)
	if err != nil {
		return nil, err
	}

	data := make([]BookDataWithCovers, 0, len(rows))
	for _, row := range rows {
		var bd BookDataWithCovers
		bd.Title = row.Title
		bd.Author = row.Author
		bd.Status = row.Status
		bd.BookID = row.BookID
		if row.CoverHash.Valid {
			bd.CoverHash = row.CoverHash.String
		} else {
			bd.CoverHash = ""
		}
		data = append(data, bd)
	}

	return data, nil
}

func (bs *BookService) FilterBySeries(ctx context.Context, accountID types.AccountID, seriesID types.SeriesID) ([]BookDataWithCovers, error) {
	rows, err := DB(ctx).GetBooksForSeries(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	data := make([]BookDataWithCovers, 0, len(rows))
	for _, row := range rows {
		var bd BookDataWithCovers
		bd.Title = row.Title
		bd.Author = row.Author
		bd.Status = row.Status
		bd.BookID = row.BookID
		if row.CoverHash.Valid {
			bd.CoverHash = row.CoverHash.String
		} else {
			bd.CoverHash = ""
		}
		data = append(data, bd)
	}

	return data, nil
}


func (bs *BookService) FilterByCollection(ctx context.Context, accountID types.AccountID, collectionID types.CollectionID) ([]BookDataWithCovers, error) {
	rows, err := DB(ctx).GetBooksForCollection(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	data := make([]BookDataWithCovers, 0, len(rows))
	for _, row := range rows {
		var bd BookDataWithCovers
		bd.Title = row.Title
		bd.Author = row.Author
		bd.Status = row.Status
		bd.BookID = row.BookID
		if row.CoverHash.Valid {
			bd.CoverHash = row.CoverHash.String
		} else {
			bd.CoverHash = ""
		}
		data = append(data, bd)
	}

	return data, nil
}
