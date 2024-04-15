package services

import (
	"context"

	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/rs/zerolog"
)

type BookService struct {
}

func NewBookService() BookService {
	return BookService{}
}

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

	return bookID, nil
}

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
