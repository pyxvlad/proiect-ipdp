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

func FixtureBook(ctx context.Context, t* testing.T) types.BookID {
	t.Helper()


	bs := services.NewBookService()

	accountID := FixtureAccount(ctx, t)

	bookID, err := bs.CreateBook(ctx, accountID, bookTitle, bookAuthor, bookStatus)
	if err != nil {
		t.Fatal(err)
	}

	return bookID
}


func TestCreateBook(t *testing.T) {
	t.Parallel()
	ctx := Context(t)

	FixtureBook(ctx, t)
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
