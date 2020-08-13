package controller

import (
	"errors"
	"net/http"

	"github.com/altescy/bookshelf/api/model"
	"github.com/julienschmidt/httprouter"
)

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.PostFormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		h.handleError(w, errors.New("username or password is empty"), http.StatusBadRequest)
		return
	}

	user, err := model.UserLogin(h.db, username, password)
	switch {
	case err == model.ErrUserNotFound:
		h.handleError(w, err, http.StatusNotFound)
	case err != nil:
		h.handleError(w, err, http.StatusInternalServerError)
	default:
		session, err := h.store.Get(r, SessionName)
		if err != nil {
			h.handleError(w, err, http.StatusInternalServerError)
			return
		}
		session.Values["user_id"] = user.ID
		if err = session.Save(r, w); err != nil {
			h.handleError(w, err, http.StatusInternalServerError)
			return
		}
		h.handleSuccess(w, user)
	}
}

func (h *Handler) Signout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, err := h.store.Get(r, SessionName)
	if err != nil {
		h.handleError(w, err, 500)
		return
	}
	session.Values["user_id"] = 0
	session.Options.MaxAge = -1
	if err = session.Save(r, w); err != nil {
		h.handleError(w, err, 500)
		return
	}
	h.handleSuccess(w, struct{}{})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := h.userByRequest(r)
	if err != nil {
		h.handleError(w, err, http.StatusUnauthorized)
		return
	}
	h.handleSuccess(w, user)
}
