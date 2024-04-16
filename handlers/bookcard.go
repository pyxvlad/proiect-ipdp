package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/pyxvlad/proiect-ipdp/templates"
)

// SampleBookCards generates 16 sample book cards, and then renders them to a page.
func SampleBookCards(w http.ResponseWriter, r *http.Request) {
	infos := make([]templates.BookCard, 0, 16)

	for i := 0; i != 16; i++ {
		var author string
		if i%2 == 0 {
			author = "no spaces"
		} else {
			author = "with spaces"
		}
		bc := templates.BookCard{
			Name:     strings.Repeat("yep"+strings.Repeat(" ", i%2), i),
			Author:   author,
			ImageURL: "https://cdn.dc5.ro/img-prod/2191826525-0.jpeg",
		}
		infos = append(infos, bc)
	}

	templates.BookCardsPage(infos).Render(context.TODO(), w)
}
