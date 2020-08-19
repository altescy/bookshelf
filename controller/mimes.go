package controller

import (
	"net/http"

	"github.com/altescy/bookshelf/model"
	"github.com/julienschmidt/httprouter"
)

func (h *Handler) GetMime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ext := "." + ps.ByName("ext")
	mime, err := model.MimeByExt(ext)
	switch {
	case err == model.ErrMimeNotFound:
		h.handleError(w, err, http.StatusNotFound)
		return
	case err != nil:
		h.handleError(w, err, http.StatusInternalServerError)
	}
	h.handleSuccess(w, mime)
}

func (h *Handler) GetMimes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	mimes := model.GetMimes()
	h.handleSuccess(w, mimes)
}
