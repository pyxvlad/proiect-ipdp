package services

import (
	"context"

	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/rs/zerolog"
)

type PublisherService struct {
}

func NewPublisherService() PublisherService {
	return PublisherService{}
}

func (ps *PublisherService) CreatePublisher(
	ctx context.Context, name string, accountID types.AccountID,
) (types.PublisherID, error) {
	publisherID, err := DB(ctx).CreatePublisher(
		ctx, database.CreatePublisherParams{
			AccountID: accountID,
			Name:      name,
		},
	)

	if err != nil {
		return 0, err
	}

	return publisherID, nil
}

type PublisherData struct {
	PublisherID types.PublisherID
	Name        string
}

func (ps *PublisherService) ListPublishers(
	ctx context.Context, accountID types.AccountID,
) ([]PublisherData, error) {
	rows, err := DB(ctx).ListPublishersForAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	publishers := make([]PublisherData, 0, len(rows))
	for _, row := range rows {
		data := PublisherData{
			PublisherID: row.PublisherID,
			Name:        row.Name,
		}
		publishers = append(publishers, data)
	}

	return publishers, nil
}

func (ps *PublisherService) RenamePublisher(
	ctx context.Context,
	accountID types.AccountID,
	publisherID types.PublisherID,
	newName string,
) error {
	log := zerolog.Ctx(ctx)
	err := DB(ctx).RenamePublisher(ctx, database.RenamePublisherParams{
		Name:        newName,
		PublisherID: publisherID,
		AccountID:   accountID,
	})

	if err != nil {
		log.Err(err).Int64("publisher_id", int64(publisherID)).Msg("couldn't rename publisher")

		return err
	}

	return err
}

func (ps *PublisherService) DeletePublisher(
	ctx context.Context,
	accountID types.AccountID,
	publisherID types.PublisherID,
) error {
	log := zerolog.Ctx(ctx)
	err := DB(ctx).DeletePublisher(ctx, database.DeletePublisherParams{
		PublisherID: publisherID,
		AccountID:   accountID,
	})

	if err != nil {
		log.Err(err).Int64("publisher_id", int64(publisherID)).Msg("couldn't delete publisher")

		return err
	}

	return err
}

func (ps *PublisherService) GetNameOfPublisher(
	ctx context.Context,
	accountID types.AccountID,
	publisherID types.PublisherID,
) (string, error) {
	return DB(ctx).GetNameOfPublisher(ctx, database.GetNameOfPublisherParams{
		AccountID: accountID,
		PublisherID: publisherID,
	})
}
