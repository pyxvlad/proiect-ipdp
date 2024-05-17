package services

import (
	"context"
	"database/sql"

	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/database/types"
)

type CollectionSeriesService struct {
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
