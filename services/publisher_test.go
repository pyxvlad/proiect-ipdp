package services_test

import (
	"context"
	"testing"

	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/pyxvlad/proiect-ipdp/services"
)

const publisherName = "Penguin"

func FixturePublisherWithSeed(ctx context.Context, t *testing.T, seed string) types.PublisherID {
	t.Helper()

	ps := services.NewPublisherService()
	accountID := FixtureAccount(ctx, t)
	publisherID, err := ps.CreatePublisher(ctx, publisherName + seed, accountID)
	if err != nil {
		t.Fatal(err)
	}

	return publisherID
}

func FixturePublisher(ctx context.Context, t *testing.T) types.PublisherID {
	t.Helper()
	return FixturePublisherWithSeed(ctx, t, "")
}
