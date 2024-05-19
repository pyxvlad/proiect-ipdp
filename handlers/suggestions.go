package handlers

import (
	"net/http"
	"strings"

	"github.com/pyxvlad/proiect-ipdp/services"
	"github.com/pyxvlad/proiect-ipdp/templates"
	"github.com/rs/zerolog"
)

func SuggestPublishers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := zerolog.Ctx(ctx)

	err := r.ParseForm()
	if err != nil {
		// TODO(ora44): Find better way to handle this
		panic(err)
	}

	publisher := r.FormValue("publisher")
	ps := services.NewPublisherService()
	data, err := ps.ListPublishers(ctx, AccountID(ctx))
	if err != nil {
		log.Err(err).Msg("while trying to list publishers")
		return
	}

	for _, pub := range data {
		if !strings.Contains(pub.Name, publisher) {
			continue
		}
		err = templates.Suggestion("publisher", uint(pub.PublisherID), pub.Name).Render(ctx, w)
		if err != nil {
			log.Err(err).Msg("while trying to display publishers")
			return
		}
	}
}

func SuggestCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := zerolog.Ctx(ctx)

	err := r.ParseForm()
	if err != nil {
		// TODO(ora44): Find better way to handle this
		panic(err)
	}

	collection := r.FormValue("collection")
	css := services.NewCollectionSeriesService()
	data, err := css.ListCollectionsForAccount(ctx, AccountID(ctx))
	if err != nil {
		log.Err(err).Msg("while trying to list collections")
		return
	}

	for _, col := range data {
		if !strings.Contains(col.Name, collection) {
			continue
		}
		err = templates.Suggestion("collection", uint(col.CollectionID), col.Name).Render(ctx, w)
		if err != nil {
			log.Err(err).Msg("while trying to display collections")
			return
		}
	}
}

func SuggestSeries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := zerolog.Ctx(ctx)

	err := r.ParseForm()
	if err != nil {
		// TODO(ora44): Find better way to handle this
		panic(err)
	}

	series := r.FormValue("series")
	css := services.NewCollectionSeriesService()
	data, err := css.ListSeriesForAccount(ctx, AccountID(ctx))
	if err != nil {
		log.Err(err).Msg("while trying to list series")
		return
	}

	for _, ser := range data {
		if !strings.Contains(ser.Name, series) {
			continue
		}
		err = templates.Suggestion("series", uint(ser.SeriesID), ser.Name).Render(ctx, w)
		if err != nil {
			log.Err(err).Msg("while trying to display series")
			return
		}
	}
}

func SuggestDuplicates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := zerolog.Ctx(ctx)

	err := r.ParseForm()
	if err != nil {
		// TODO(ora44): Find better way to handle this
		panic(err)
	}

	duplicates := r.FormValue("duplicate")
	bs := services.NewBookService(services.CoverPath(ctx))
	data, err := bs.ListBooksForAccount(ctx, AccountID(ctx))
	if err != nil {
		log.Err(err).Msg("while trying to list duplicates")
		return
	}

	for _, dup := range data {
		if !strings.Contains(dup.Title, duplicates) {
			continue
		}
		err = templates.Suggestion("duplicate", uint(dup.BookID), dup.Title+" - "+dup.Author).Render(ctx, w)
		if err != nil {
			log.Err(err).Msg("while trying to display duplicates")
			return
		}
	}
}
