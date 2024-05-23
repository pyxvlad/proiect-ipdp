package handlers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/pyxvlad/proiect-ipdp/templates"
	"github.com/rs/zerolog"
)

// SampleBookCards generates 16 sample book cards, and then renders them to a page.
/*func SampleBookCards(w http.ResponseWriter, r *http.Request) {
	infos := make([]templates.BookCard, 0, 24)

	for i := 0; i != 24; i++ {
		var author string
		if i%2 == 0 {
			author = "no spaces"
		} else {
			author = "with spaces"
		}
		bc := templates.BookCard{
			Title:    strings.Repeat("yep"+strings.Repeat(" ", i%2), i),
			Author:   author,
			ImageURL: "https://cdn.dc5.ro/img-prod/2191826525-0.jpeg",
			Status:   types.StatusToBeRead,
		}
		infos = append(infos, bc)
	}

	templates.BookCardsPage(infos).Render(context.TODO(), w)
}*/

func ViewLibrary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := zerolog.Ctx(ctx)
	bs := services.NewBookService(services.CoverPath(ctx))
	infos, err := bs.ListBooksWithCoversForAccount(ctx, AccountID(ctx))
	if err != nil {
		log.Err(err).Msg("while trying to list books")
		return
	}
	templates.BookCardsPage(infos).Render(context.TODO(), w)
}

func PreviewCard(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		// TODO(ora44): Find better way to handle this
		panic(err)
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	statusRaw := r.FormValue("status")
	file, header, err := r.FormFile("cover")
	_ = header
	var dataURL = ""

	if err == nil || !errors.Is(err, http.ErrMissingFile) {
		if err != nil {
			log.Err(err).Msg("Could not upload image")
			return
		}
		data, err := io.ReadAll(file)
		if err != nil {
			log.Err(err).Msg("Could not read all image")
			return
		}
		encoded64 := base64.StdEncoding.EncodeToString(data)
		dataURL = fmt.Sprintf("data:image/png;base64,%s", encoded64)
	}

	err = templates.BookCardPreview(templates.BookCard{
		Title:    title,
		Author:   author,
		ImageURL: dataURL,
		Status:   types.Status(statusRaw),
	}).Render(r.Context(), w)
	if err != nil {
		log.Err(err).Msg("Trouble updating the preview")
	}
}

func ViewDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := zerolog.Ctx(ctx)

	bookIDString := chi.URLParam(r, "bookID")
	parsedID, err := strconv.ParseInt(bookIDString, 10, 64)
	if err != nil {
		parsedID = 0
	}

	bookID := types.BookID(parsedID)

	bs := services.NewBookService(services.CoverPath(ctx))

	bookData, err := bs.GetAllDataForBook(ctx, bookID, AccountID(ctx))
	if err != nil {
		log.Err(err).Msg("while getting book data")
		return
	}

	err = templates.DetailsPage(bookData).Render(ctx, w)
	if err != nil {
		log.Err(err).Msg("while rendering details page")
	}
}
