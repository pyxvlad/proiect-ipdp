package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io"
	"os"
	"path"

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
		AccountID:  account,
		Title:      title,
		Author:     author,
		ProgressID: progressID,
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
