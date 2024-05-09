package services

import (
	"context"
	"database/sql"

	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/rs/zerolog"
)

type BookService struct {
}

func NewBookService() BookService {
	return BookService{}
}

// CreateBook creates a new book.
func (b *BookService) CreateBook(
	ctx context.Context,
	account types.AccountID,
	title string,
	author string,
	status types.Status,
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
