package services

import (
	"context"
	"database/sql"

	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/database/types"
)

type CollectionSeriesService struct {
}

func NewCollectionSeriesService() CollectionSeriesService {
	return CollectionSeriesService{}
}

func (css *CollectionSeriesService) CreateSeries(ctx context.Context, accountID types.AccountID, name string) (types.SeriesID, error) {
	return DB(ctx).CreateSeries(ctx, database.CreateSeriesParams{
		Name:      name,
		AccountID: accountID,
	})
}

func (css *CollectionSeriesService) CreateCollection(ctx context.Context, accountID types.AccountID, name string) (types.CollectionID, error) {
	return DB(ctx).CreateCollection(ctx, database.CreateCollectionParams{
		Name:      name,
		AccountID: accountID,
	})
}

func (css *CollectionSeriesService) AddBookToSeries(ctx context.Context, bookID types.BookID, seriesID types.SeriesID, volume uint) error {
	return DB(ctx).AddBookToSeries(ctx, database.AddBookToSeriesParams{
		SeriesID: seriesID,
		BookID:   bookID,
		Volume:   volume,
	})
}

// Adds a book to a collection. Using book_number=0 will leave the book_number as null/nil.
func (css *CollectionSeriesService) AddBookToCollection(ctx context.Context, bookID types.BookID, collectionID types.CollectionID, book_number uint) error {
	return DB(ctx).AddBookToCollection(ctx, database.AddBookToCollectionParams{
		CollectionID: collectionID,
		BookID:       bookID,
		BookNumber: sql.NullInt64{
			Int64: int64(book_number),
			Valid: book_number != 0,
		},
	})
}

func (css *CollectionSeriesService) RemoveBookFromSeries(ctx context.Context, bookID types.BookID, seriesID types.SeriesID) error {
	return DB(ctx).RemoveBookFromSeries(ctx, database.RemoveBookFromSeriesParams{
		SeriesID: seriesID,
		BookID:   bookID,
	})
}

func (css *CollectionSeriesService) RemoveBookFromCollection(ctx context.Context, bookID types.BookID, collectionID types.CollectionID) error {
	return DB(ctx).RemoveBookFromCollection(ctx, database.RemoveBookFromCollectionParams{
		CollectionID: collectionID,
		BookID:       bookID,
	})
}

func (css *CollectionSeriesService) DeleteCollection(ctx context.Context, collectionID types.CollectionID, accountID types.AccountID) error {
	return DB(ctx).DeleteCollection(ctx, database.DeleteCollectionParams{
		CollectionID: collectionID,
		AccountID:    accountID,
	})
}

func (css *CollectionSeriesService) DeleteSeries(ctx context.Context, seriesID types.SeriesID, accountID types.AccountID) error {
	return DB(ctx).DeleteSeries(ctx, database.DeleteSeriesParams{
		SeriesID:  seriesID,
		AccountID: accountID,
	})
}

type SeriesData struct {
	SeriesID types.SeriesID
	Name     string
}

func (css *CollectionSeriesService) ListSeriesForAccount(ctx context.Context, accountID types.AccountID) ([]SeriesData, error) {
	row, err := DB(ctx).ListSeriesForAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	data := make([]SeriesData, 0, len(row))
	for _, v := range row {
		data = append(data, SeriesData{
			SeriesID: v.SeriesID,
			Name:     v.Name,
		})
	}
	return data, nil
}

type CollectionData struct {
	CollectionID types.CollectionID
	Name         string
}

func (css *CollectionSeriesService) ListCollectionsForAccount(ctx context.Context, accountID types.AccountID) ([]CollectionData, error) {
	row, err := DB(ctx).ListCollectionsForAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	data := make([]CollectionData, 0, len(row))
	for _, v := range row {
		data = append(data, CollectionData{
			CollectionID: v.CollectionID,
			Name:         v.Name,
		})
	}
	return data, nil
}

func (css *CollectionSeriesService) RenameSeries(ctx context.Context, accountID types.AccountID, seriesID types.SeriesID, name string) error {
	return DB(ctx).RenameSeries(ctx, database.RenameSeriesParams{
		Name:      name,
		SeriesID:  seriesID,
		AccountID: accountID,
	})
}

func (css *CollectionSeriesService) RenameCollection(ctx context.Context, accountID types.AccountID, collectionID types.CollectionID, name string) error {
	return DB(ctx).RenameCollection(ctx, database.RenameCollectionParams{
		Name:         name,
		CollectionID: collectionID,
		AccountID:    accountID,
	})
}

func (css *CollectionSeriesService) SetVolumeInSeries(ctx context.Context, volume uint, bookID types.BookID) error {
	return DB(ctx).ChangeBookVolumeInSeries(ctx, database.ChangeBookVolumeInSeriesParams{
		Volume: volume,
		BookID: bookID,
	})
}

func (css *CollectionSeriesService) SetNumberInCollection(ctx context.Context, number uint, bookID types.BookID) error {
	return DB(ctx).ChangeBookNumberInCollection(ctx, database.ChangeBookNumberInCollectionParams{
		BookNumber: sql.NullInt64{
			Int64: int64(number),
			Valid: number != 0,
		},
		BookID: bookID,
	})
}

type BookSeriesData struct {
	SeriesID types.SeriesID
	Name     string
	Volume   uint
}

func (css *CollectionSeriesService) GetSeriesForBook(
	ctx context.Context, bookID types.BookID, accountID types.AccountID,
) (BookSeriesData, error) {
	row, err := DB(ctx).GetSeriesForBook(ctx, database.GetSeriesForBookParams{
		AccountID: accountID,
		BookID:    bookID,
	})
	if err != nil {
		return BookSeriesData{}, err
	}

	return BookSeriesData{
		SeriesID: row.SeriesID,
		Name:     row.Name,
		Volume:   row.Volume,
	}, nil
}

type BookCollectionData struct {
	CollectionID types.CollectionID
	Name         string
	Number       uint
}

func (css *CollectionSeriesService) GetCollectionForBook(
	ctx context.Context, bookID types.BookID, accountID types.AccountID,
) (BookCollectionData, error) {
	row, err := DB(ctx).GetCollectionForBook(ctx, database.GetCollectionForBookParams{
		AccountID: accountID,
		BookID:    bookID,
	})
	if err != nil {
		return BookCollectionData{}, err
	}

	var number uint
	if row.BookNumber.Valid {
		number = uint(row.BookNumber.Int64)
	} else {
		number = 0
	}

	return BookCollectionData{
		CollectionID: row.CollectionID,
		Name:         row.Name,
		Number:       number,
	}, nil
}
