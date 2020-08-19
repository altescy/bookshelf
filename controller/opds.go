package controller

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/altescy/bookshelf/model"
	"github.com/altescy/bookshelf/opds"
	"github.com/julienschmidt/httprouter"
)

const opdsTitle = "Bookshelf - OPDS"

func (h *Handler) GetOPDSFeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	books, err := model.GetBooks(h.db)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	opdsURL := "http://" + r.Host + r.URL.Path
	href := r.URL.Path

	for i, book := range *books {
		for j, file := range book.Files {
			alias, _ := model.GetMimeAlias(file.MimeType)
			(*books)[i].Files[j].Link = fmt.Sprintf("/api/book/%d/file/%s", book.ID, alias)
		}
	}

	entries := model.EntriesFromBooks(books)
	feed := opds.BuildFeed(opdsURL, opdsTitle, href, entries)

	enc := xml.NewEncoder(w)

	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, xml.Header)
	err = enc.Encode(&feed)
	if err != nil {
		log.Printf("[error] cannot encode xml feed: %v", err)
		if err != io.ErrClosedPipe {
			h.handleError(w, err, http.StatusInternalServerError)
			return
		}
	}
}
