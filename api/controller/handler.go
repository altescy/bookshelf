package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) CommonMiddleware(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				h.handleError(w, err, 400)
				return
			}
		}
		f.ServeHTTP(w, r)
	})
}

func (h *Handler) handleSuccess(w http.ResponseWriter, data interface{}) {
	// DEV: enableCors(&w)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[WARN] write response json failed. %s", err)
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error, code int) {
	// DEV: enableCors(&w)
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

func (h *Handler) txScope(f func(*sql.Tx) error) (err error) {
	var tx *sql.Tx
	tx, err = h.db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin transaction failed")
	}
	defer func() {
		if e := recover(); e != nil {
			tx.Rollback()
			err = errors.Errorf("panic in transaction: %s", e)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = f(tx)
	return
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hello, world!")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
