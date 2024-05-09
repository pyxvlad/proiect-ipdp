package services_test

import (
	"context"
	"testing"

	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/pyxvlad/proiect-ipdp/services"
)

const (
	bookTitle = "Crimă și pedeapsă"
	bookAuthor = "Dostoievski"
	bookStatus = types.StatusRead
)

func FixtureBookWithSeed(ctx context.Context, t* testing.T, seed string) types.BookID {
	t.Helper()

	bs := services.NewBookService()

	accountID := FixtureAccount(ctx, t)

	bookID, err := bs.CreateBook(ctx, accountID, bookTitle + " " + seed, bookAuthor, bookStatus)
	if err != nil {
		t.Fatal(err)
	}

	return bookID
}

func FixtureBook(ctx context.Context, t *testing.T) types.BookID{
	t.Helper()

	return FixtureBookWithSeed(ctx, t, "")
}

func TestCreateBook(t *testing.T) {
	t.Parallel()
	ctx := Context(t)

	bookID := FixtureBook(ctx, t)

	// This is used just to check that the book _exists_ in the database.

	_, err := services.DB(ctx).GetDuplicateOfBook(ctx, bookID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetBookStatus(t *testing.T) {
	t.Parallel()
	ctx := Context(t)

	bs := services.NewBookService()
	bookID := FixtureBook(ctx, t)

	err := bs.SetBookStatus(ctx, bookID, types.StatusInProgress)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMarkBookAsDuplicate(t *testing.T) {
	t.Parallel()

	ctx := Context(t)

	bs := services.NewBookService()

	firstBookID := FixtureBookWithSeed(ctx, t, "first")
	duplicateBookID := FixtureBookWithSeed(ctx, t, "duplicate")

	independentBookID := FixtureBookWithSeed(ctx, t, "independent")

	err := bs.MarkBookAsDuplicate(ctx, firstBookID, duplicateBookID)
	if err != nil {
		t.Fatal(err)
	}

	db := services.DB(ctx)
	firstDuplicateID, err := db.GetDuplicateOfBook(ctx, firstBookID)
	if err != nil {
		t.Fatal(err)
	}
	
	duplicatedDuplicateID, err := db.GetDuplicateOfBook(ctx, duplicateBookID)
	if err != nil {
		t.Fatal(err)
	}

	if firstDuplicateID != duplicatedDuplicateID {
		t.Fatal("duplicate IDs do not match")
	}

	independentDuplicateID, err := db.GetDuplicateOfBook(ctx, independentBookID)
	if err != nil {
		t.Fatal(err)
	}

	if firstDuplicateID == independentDuplicateID {
		t.Fatal("independent book was marked as duplicate")
	}


}
