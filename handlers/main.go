package handlers

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/pyxvlad/proiect-ipdp/database/types"
	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/pyxvlad/proiect-ipdp/templates"
	"github.com/rs/zerolog"
)

type HelloHandler struct {
	log zerolog.Logger
}

func NewHelloHandler(log zerolog.Logger) http.Handler {
	handler := new(HelloHandler)

	handler.log = log

	return handler
}

// ServeHTTP implements http.Handler.
func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := templates.HelloPage().Render(r.Context(), w)
	h.log.Err(err).Send()
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())
	err := templates.SignUpPage().Render(r.Context(), w)
	log.Err(err).Send()
}

func SignUpAttempt(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())

	err := r.ParseForm()
	if err != nil {
		// TODO(ora44): Find better way to handle this
		panic(err)
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	_, err = mail.ParseAddress(email)
	if err != nil {
		err = templates.SignUpForm("Invalid email address, retry.").Render(r.Context(), w)
		log.Err(err).Send()
		return
	}

	if password != confirmPassword {
		err = templates.SignUpForm("Passwords don't match, retry.").Render(r.Context(), w)
		log.Err(err).Send()
		return
	}

	as := services.NewAccountService()
	as.CreateAccountWithEmail(r.Context(), services.AccountData{
		Email:    email,
		Password: password,
	})

	w.Header().Add("HX-Location", "/hello")

	w.WriteHeader(http.StatusCreated)
}

func LogInPage(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())
	err := templates.LogInPage().Render(r.Context(), w)
	log.Err(err).Send()
}

func LogInAttempt(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	accountService := services.NewAccountService()
	accountID, err := accountService.Login(r.Context(), services.AccountData{
		Email:    email,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	token, err := accountService.CreateSession(r.Context(), accountID)
	if err != nil {
		panic(err)
	}
	cookie := http.Cookie{
		Name:   "token",
		Value:  token,
		MaxAge: 0,
		Path:   "/",
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/hello")
	_ = log
}

func IDFromForm[T ~int64](r *http.Request, key string) T {
	parsedID, err := strconv.ParseInt(r.FormValue(key), 10, 64)
	if err != nil {
		parsedID = 0
	}
	return T(parsedID)
}

func AddBookPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := zerolog.Ctx(ctx)
	log.Debug().Str("http_method", r.Method).Send()
	if r.Method == http.MethodGet {
		err := templates.AddBookPage().Render(r.Context(), w)
		log.Err(err).Send()
		return
	}
	if r.Method != http.MethodPost {
		panic("Wrong method")
	}

	err := r.ParseMultipartForm(10 << 20)
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
	duplicateID := IDFromForm[types.BookID](r, "book_id")

	// TODO: use header
	_ = header

	bs := services.NewBookService(services.CoverPath(ctx))

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

	bookID, err := bs.CreateBook(ctx, AccountID(ctx), title, author, types.Status(statusRaw), publisherID)
	if err != nil {
		log.Err(err).Msg("while trying to create book")
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

}

func Menu(w http.ResponseWriter, r *http.Request) {
	log := zerolog.Ctx(r.Context())
	err := templates.Menu().Render(r.Context(), w)
	log.Err(err).Send()
}
