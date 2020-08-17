package controller

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/altescy/bookshelf/api/model"
	"github.com/altescy/bookshelf/api/opds"
	"github.com/julienschmidt/httprouter"
)

const opdsTitle = "bookshelf - OPDS"

func (h *Handler) GetOPDSFeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	books, err := model.GetBooks(h.db)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	// assume that opdsURL is http://[host][basePath]/opds
	opdsURL := "http://" + r.Host + r.URL.Path
	basePath := r.URL.Path[:len(r.URL.Path)-5]
	href := r.URL.Path

	for i, book := range *books {
		for j, file := range book.Files {
			alias, _ := model.GetMimeAlias(file.MimeType)
			(*books)[i].Files[j].Link = fmt.Sprintf("%s/book/%d/file/%s", basePath, book.ID, alias)
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
