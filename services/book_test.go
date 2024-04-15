package services_test

import (
	"testing"

	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/pyxvlad/proiect-ipdp/services"
)

const (
	bookTitle = "Crimă și pedeapsă"
	bookAuthor = "Dostoievski"
	bookStatus = types.StatusRead
)


func TestCreateBook(t *testing.T) {
	ctx := Context(t)

	bs := services.NewBookService()

	accountID := FixtureAccount(ctx, t)

	_, err := bs.CreateBook(ctx, accountID, bookTitle, bookAuthor, bookStatus)
	if err != nil {
		t.Fatal(err)
	}
}
