package services_test

import (
	"bytes"
	"context"
	_ "embed"
	"os"
	"path"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/pyxvlad/proiect-ipdp/services"
)

const (
	bookTitle  = "Crimă și pedeapsă"
	bookAuthor = "Dostoievski"
	bookStatus = types.StatusRead
)

func FixtureBookWithSeed(ctx context.Context, t *testing.T, seed string) types.BookID {
	t.Helper()

	bs := FixtureBookService(ctx, t)

	accountID := FixtureAccount(ctx, t)

	bookID, err := bs.CreateBook(ctx, accountID, bookTitle+" "+seed, bookAuthor, bookStatus)
	if err != nil {
		t.Fatal(err)
	}

	return bookID
}

func FixtureBook(ctx context.Context, t *testing.T) types.BookID {
	t.Helper()

	return FixtureBookWithSeed(ctx, t, "")
}

func FixtureBookService(ctx context.Context, t *testing.T) services.BookService {
	t.Helper()

	mapper := func(r rune) rune {
		if r < utf8.RuneSelf {
			const allowed = "!#$%&()+,-.=@^_{}~ "
			if '0' <= r && r <= '9' ||
				'a' <= r && r <= 'z' ||
				'A' <= r && r <= 'Z' {
				return r
			}
			if strings.ContainsRune(allowed, r) {
				return r
			}
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}
	pattern := strings.Map(mapper, t.Name())

	pattern = pattern + "-images"

	// TODO: rename databasePath, it might get confusing
	imagePath := path.Join(databasePath, pattern)

	err := os.MkdirAll(imagePath, 0777)
	if err != nil {
		t.Fatal(err)
	}

	return services.NewBookService(imagePath)
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

	bs := FixtureBookService(ctx, t)
	bookID := FixtureBook(ctx, t)

	err := bs.SetBookStatus(ctx, bookID, types.StatusInProgress)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMarkBookAsDuplicate(t *testing.T) {
	t.Parallel()

	ctx := Context(t)

	bs := FixtureBookService(ctx, t)

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

func TestListBooks(t *testing.T) {
	t.Parallel()
	ctx := Context(t)

	bs := FixtureBookService(ctx, t)
	bookIDA := FixtureBookWithSeed(ctx, t, "A")
	bookIDB := FixtureBookWithSeed(ctx, t, "B")

	accountID := FixtureAccount(ctx, t)

	books, err := bs.ListBooksForAccount(ctx, accountID)
	if err != nil {
		t.Fatal(err)
	}

	if len(books) != 2 {
		t.Fatal("Expected 2 books, got", len(books))
	}

	for _, book := range books {
		// TODO(vlad): use table testing for this maybe?
		if book.BookID == bookIDA && (book.Title != (bookTitle + " " + "A")) {
			t.Fatal("wrong book title on book A: ", book.Title)
		}

		if book.BookID == bookIDB && (book.Title != (bookTitle + " " + "B")) {
			t.Fatal("wrong book title on book B: ", book.Title)
		}
	}

}

//go:embed testdata/ps2n2c16.png
var samplePNG []byte

func TestSetBookCover(t *testing.T) {
	t.Parallel()
	ctx := Context(t)

	bs := FixtureBookService(ctx, t)

	os.MkdirAll("/tmp/ipdp-img", 0777)

	bookID := FixtureBook(ctx, t)

	err := bs.SetBookCover(ctx, bookID, bytes.NewReader(samplePNG))
	if err != nil {
		t.Fatal(err)
	}


	rows, err := services.DB(ctx).GetBooksWithCoversAndStatuses(ctx, FixtureAccount(ctx, t))
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 1 {
		t.Fatal("expected one book")
	}

	if rows[0].BookID != bookID {
		t.Fatal("got the wrong book")
	}

	if !rows[0].CoverHash.Valid {
		t.Fatal("expected the book to have the cover hash set")
	}

	data, err := os.ReadFile(path.Join(bs.ImagePath, rows[0].CoverHash.String))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(samplePNG, data) {
		t.Fatal("the sample file and the saved one differ")
	}
	
}
