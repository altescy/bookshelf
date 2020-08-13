package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/altescy/bookshelf/api/model"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

type key int

const (
	keyUserID key = iota
)

const (
	SessionName = "bookshelf_session"
)

type Handler struct {
	db    *gorm.DB
	store sessions.Store
}

func NewHandler(db *gorm.DB, store sessions.Store) *Handler {
	return &Handler{
		db:    db,
		store: store,
	}
}

func (h *Handler) CommonMiddleware(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				h.handleError(w, err, http.StatusBadRequest)
				return
			}
		}
		session, err := h.store.Get(r, SessionName)
		if err != nil {
			h.handleError(w, err, http.StatusInternalServerError)
			return
		}
		if _userID, ok := session.Values["user_id"]; ok {
			userID := _userID.(uint)
			user, err := model.GetUserByID(h.db, userID)
			switch {
			case err == sql.ErrNoRows:
				session.Values["user_id"] = 0
				session.Options = &sessions.Options{MaxAge: -1}
				if err = session.Save(r, w); err != nil {
					h.handleError(w, err, http.StatusInternalServerError)
					return
				}
				h.handleError(w, errors.New("session disconnected"), http.StatusNotFound)
				return
			case err != nil:
				h.handleError(w, err, http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), keyUserID, user.ID)
			f.ServeHTTP(w, r.WithContext(ctx))
		} else {
			f.ServeHTTP(w, r)
		}
	})
}

func (h *Handler) userByRequest(r *http.Request) (*model.User, error) {
	v := r.Context().Value(keyUserID)
	if id, ok := v.(uint); ok {
		return model.GetUserByID(h.db, id)
	}
	return nil, errors.New("Not authenticated")
}

func (h *Handler) handleSuccess(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[WARN] write response json failed. %s", err)
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	log.Printf("[WARN] err: %s", err.Error())
	data := map[string]interface{}{
		"code": code,
		"err":  err.Error(),
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[WARN] write error response json failed. %s", err)
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, world!")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
