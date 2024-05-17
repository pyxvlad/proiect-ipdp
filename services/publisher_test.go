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
	publisherID, err := ps.CreatePublisher(ctx, publisherName+seed, accountID)
	if err != nil {
		t.Fatal(err)
	}

	return publisherID
}

func FixturePublisher(ctx context.Context, t *testing.T) types.PublisherID {
	t.Helper()
	return FixturePublisherWithSeed(ctx, t, "")
}

func TestCreatePublisher(t *testing.T) {
	ctx := Context(t)
	publisherID := FixturePublisher(ctx, t)

	ps := services.NewPublisherService()

	publishersData, err := ps.ListPublishers(ctx, FixtureAccount(ctx, t))
	if err != nil {
		t.Fatal(err)
	}

	if len(publishersData) != 1 {
		t.Fatal("expected to have one publisher")
	}

	if publishersData[0].PublisherID != publisherID {
		t.Fatalf(
			"publisher IDs don't match, got:%d expected:%d",
			publishersData[0].PublisherID, publisherID,
		)
	}

	if publishersData[0].Name != publisherName {
		t.Fatalf(
			"publisher names don't match, got:%s expected:%s",
			publishersData[0].Name, publisherName,
		)
	}
}

func TestRenamePublisher(t *testing.T) {
	const newName = "Nemira"
	ctx := Context(t)

	ps := services.NewPublisherService()

	publisherID := FixturePublisher(ctx, t)

	accountID := FixtureAccount(ctx, t)
	err := ps.RenamePublisher(ctx, accountID, publisherID, newName)
	if err != nil {
		t.Fatal(err)
	}

	publishers, err := ps.ListPublishers(ctx, accountID)
	if err != nil {
		t.Fatal(err)
	}

	if len(publishers) != 1 ||
		publishers[0].PublisherID != publisherID ||
		publishers[0].Name != newName {

		t.Fatalf(
			"tried to rename %d to %s, got: %#v",
			publisherID, newName, publishers,
		)
	}
}

func TestDeletePublisher(t *testing.T) {
	ctx := Context(t)

	ps := services.NewPublisherService()

	publisherID := FixturePublisher(ctx, t)

	accountID := FixtureAccount(ctx, t)
	err := ps.DeletePublisher(ctx, accountID, publisherID)
	if err != nil {
		t.Fatal(err)
	}

	publishers, err := ps.ListPublishers(ctx, accountID)
	if err != nil {
		t.Fatal(err)
	}

	if len(publishers) != 0 {
		t.Fatalf(
			"tried to delete %d, got: %#v",
			publisherID, publishers,
		)
	}
}
