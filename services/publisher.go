package services

import (
	"context"

	"github.com/pyxvlad/proiect-ipdp/database"
	"github.com/pyxvlad/proiect-ipdp/database/types"
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
			Name: name,
		},
	)

	if err != nil {
		return 0, err
	}

	return publisherID, nil
}

type PublisherData struct {
	PublisherID types.PublisherID
	Name string

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
		data := PublisherData {
			PublisherID: row.PublisherID,
			Name: row.Name,
		}
		publishers = append(publishers, data)
	}

	return publishers, nil
}

