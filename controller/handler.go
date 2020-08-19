package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/altescy/bookshelf/storage"
	"github.com/jinzhu/gorm"
)

type key int

const (
	keyUserID key = iota
)

type Handler struct {
	db         *gorm.DB
	storage    storage.Storage
	enableCors bool
}

func NewHandler(db *gorm.DB, storage storage.Storage, enableCors bool) *Handler {
	return &Handler{
		db:         db,
		storage:    storage,
		enableCors: enableCors,
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
		if h.enableCors {
			enableCors(&w)
		}
		f.ServeHTTP(w, r)
	})
}

func (h *Handler) handleSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[WARN] write response json failed. %s", err)
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	log.Printf("[WARN] err: %s", err.Error())
	data := map[string]interface{}{
		"code": code,
		"err":  err.Error(),
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[WARN] write error response json failed. %s", err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
