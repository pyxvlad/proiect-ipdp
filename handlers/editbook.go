package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/pyxvlad/proiect-ipdp/templates"
	"github.com/rs/zerolog"
)

func EditBookPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := zerolog.Ctx(ctx)
	bs := services.NewBookService(services.CoverPath(ctx))
	bookIDString := chi.URLParam(r, "bookID")
	parsedID, err := strconv.ParseInt(bookIDString, 10, 64)
	if err != nil {
		parsedID = 0
	}

	bookID := types.BookID(parsedID)
	log.Debug().Str("http_method", r.Method).Send()
	if r.Method == http.MethodGet {
		bookData, err := bs.GetAllDataForBook(ctx, bookID, AccountID(ctx))
		if err != nil {
			log.Err(err).Msg("while trying to get all data for book")
			return
		}
		err = templates.EditBookPage(bookData).Render(r.Context(), w)
		log.Err(err).Send()
		return
	}
	if r.Method != http.MethodPost {
		panic("Wrong method")
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		// TODO(ora44): Find better way to handle this
		panic(err)
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	statusRaw := r.FormValue("status")

	file, header, coverErr := r.FormFile("cover")
	publisher := r.FormValue("publisher")
	collection := r.FormValue("collection")
	series := r.FormValue("series")
	publisherID := IDFromForm[types.PublisherID](r, "publisher_id")
	collectionID := IDFromForm[types.CollectionID](r, "collection_id")
	seriesID := IDFromForm[types.SeriesID](r, "series_id")
	duplicateID := IDFromForm[types.BookID](r, "duplicate_id")

	// TODO: use header
	_ = header

	err = bs.SetBookTitle(ctx, AccountID(ctx), bookID, title)
	if err != nil {
		log.Err(err).Msg("while trying to set book title")
		return
	}

	err = bs.SetBookAuthor(ctx, AccountID(ctx), bookID, author)
	if err != nil {
		log.Err(err).Msg("while trying to set author name")
		return
	}

	err = bs.SetBookStatus(ctx, bookID, types.Status(statusRaw))
	if err != nil {
		log.Err(err).Msg("while trying to set book status")
		return
	}

	css := services.NewCollectionSeriesService()
	ps := services.NewPublisherService()
	if publisherID == types.InvalidPublisherID && publisher != "" {
		publisherID, err = ps.CreatePublisher(ctx, publisher, AccountID(ctx))
		if err != nil {
			log.Err(err).Msg("while trying to create publisher")
			return
		}
	} else if publisherID == types.InvalidPublisherID && publisher == "" {
		log.Error().Msg("while ceva")
		// TODO: show it to the user
		return
	}

	err = bs.SetBookPublisher(ctx, bookID, publisherID)
	if err != nil {
		log.Err(err).Msg("while trying to set publisher")
		return
	}

	if coverErr == nil || !errors.Is(coverErr, http.ErrMissingFile) {
		if coverErr != nil {
			log.Err(coverErr).Msg("Could not upload image")
			return
		}

		err = bs.SetBookCover(ctx, bookID, file)
		if err != nil {
			log.Err(err).Msg("while trying to set cover image")
			return
		}
	}

	collectionNumber := IDFromForm[int64](r, "collection-numeric")

	collectionData, err := css.GetCollectionForBook(ctx, bookID, AccountID(ctx))
	if err != nil {
		log.Err(err).Msg("while getting collection ID")
	}
	if collectionID == collectionData.CollectionID {
		err := css.RemoveBookFromCollection(ctx, bookID, collectionID)
		if err != nil {
			log.Err(err).Msg("while trying to remove book from collection")
			return
		}
	}

	if collectionID != types.InvalidCollectionID {
		err = css.AddBookToCollection(ctx, bookID, collectionID, uint(collectionNumber))
		if err != nil {
			log.Err(err).Msg("while trying to add book to collection")
			return
		}
	} else if collection != "" {
		collectionID, err = css.CreateCollection(ctx, AccountID(ctx), collection)
		if err != nil {
			log.Err(err).Msg("while trying to create collection")
			return
		}
		err = css.AddBookToCollection(ctx, bookID, collectionID, uint(collectionNumber))
		if err != nil {
			log.Err(err).Msg("while trying to add book to collection")
			return
		}
	}

	seriesNumber := IDFromForm[int64](r, "series-numeric")

	seriesData, err := css.GetSeriesForBook(ctx, bookID, AccountID(ctx))
	if err != nil {
		log.Err(err).Msg("while getting series ID")
	}
	if seriesID == seriesData.SeriesID {
		err := css.RemoveBookFromSeries(ctx, bookID, seriesID)
		if err != nil {
			log.Err(err).Msg("while trying to remove book from series")
			return
		}
	}

	if seriesID != types.InvalidSeriesID {
		err = css.AddBookToSeries(ctx, bookID, seriesID, uint(seriesNumber))
		if err != nil {
			log.Err(err).Msg("while trying to add book to series")
			return
		}
	} else if series != "" {
		seriesID, err = css.CreateSeries(ctx, AccountID(ctx), series)
		if err != nil {
			log.Err(err).Msg("while trying to create series")
			return
		}
		err = css.AddBookToSeries(ctx, bookID, seriesID, uint(seriesNumber))
		if err != nil {
			log.Err(err).Msg("while trying to add book to series")
			return
		}
	}

	if duplicateID != types.InvalidBookID {
		err = bs.MarkBookAsDuplicate(ctx, types.BookID(bookID), duplicateID)
		if err != nil {
			log.Err(err).Msg("while trying to mark book as duplicate")
			return
		}
	}
	w.Header().Add("Location", fmt.Sprintf("/books/%d/details", bookID))
	w.WriteHeader(http.StatusSeeOther)

}
