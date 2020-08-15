package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/altescy/bookshelf/api/model"
	"github.com/julienschmidt/httprouter"
)

const pubdateLayout = "2006-01-02"

// AddBook add a new book into database
func (h *Handler) AddBook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	book := model.Book{
		ISBN:        r.FormValue("ISBN"),
		Title:       r.FormValue("Title"),
		Author:      r.FormValue("Author"),
		Description: r.FormValue("Description"),
		CoverURL:    r.FormValue("CoverUrl"),
		Publisher:   r.FormValue("Publisher"),
	}

	pubdateString := r.FormValue("PubDate")
	if pubdateString != "" {
		pubdate, err := time.Parse(pubdateLayout, pubdateString)
		if err != nil {
			h.handleError(w, errors.New("invalid pubdate format"), http.StatusBadRequest)
			return
		}
		*book.PubDate = pubdate
	}

	if err := model.AddBook(h.db, &book); err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	h.handleSuccess(w, book)
}

//DeleteBook delete specified book from DB
func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookidString := ps.ByName("bookid")
	bookID, err := strconv.ParseUint(bookidString, 10, 64)
	if err != nil {
		h.handleError(w, errors.New("invalid bookid"), http.StatusBadRequest)
		return
	}
	if err := model.DeleteBook(h.db, &model.Book{ID: bookID}); err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}
	h.handleSuccess(w, "successfully deleted")
}

// GetBook return a book having a specified bookid
func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookidString := ps.ByName("bookid")
	bookID, err := strconv.ParseUint(bookidString, 10, 64)
	if err != nil {
		h.handleError(w, errors.New("invalid bookid"), http.StatusBadRequest)
		return
	}

	book, err := model.GetBookByID(h.db, bookID)
	switch {
	case err == model.ErrBookNotFound:
		h.handleError(w, err, http.StatusNotFound)
		return
	case err != nil:
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	h.handleSuccess(w, book)
}

// GetBooks returns list of books where next <= bookid < next+count
func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query()

	nextString := q.Get("next")
	next, err := strconv.ParseUint(nextString, 10, 64)
	if err != nil && nextString != "" {
		h.handleError(w, errors.New("invalid next value"), http.StatusBadRequest)
		return
	}

	countString := q.Get("count")
	count, err := strconv.ParseUint(countString, 10, 64)
	if err != nil && countString != "" {
		h.handleError(w, errors.New("invalid count value"), http.StatusBadRequest)
		return
	}

	respond := func(books *[]model.Book, err error) {
		switch {
		case err == model.ErrBookNotFound:
			h.handleError(w, err, http.StatusNotFound)
			return
		case err != nil:
			h.handleError(w, err, http.StatusInternalServerError)
			return
		}
		h.handleSuccess(w, books)
	}

	switch {
	case countString != "" && nextString != "":
		respond(model.GetBooksWithNextCount(h.db, next, count))
		return
	case nextString != "" && countString == "":
		respond(model.GetBooksWithNext(h.db, next))
		return
	case countString != "" && nextString == "":
		respond(model.GetBooksWithCount(h.db, count))
		return
	default:
		respond(model.GetBooks(h.db))
		return
	}
}

// UpdateBook update book properties
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookidString := ps.ByName("bookid")
	bookID, err := strconv.ParseUint(bookidString, 10, 64)
	if err != nil {
		h.handleError(w, errors.New("invalid bookid"), http.StatusBadRequest)
	}

	book, err := model.GetBookByID(h.db, bookID)
	switch {
	case err == model.ErrBookNotFound:
		h.handleError(w, err, http.StatusNotFound)
		return
	case err != nil:
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	updateString := func(field string, value *string) {
		newValue := r.FormValue(field)
		if newValue != "" {
			*value = newValue
		}
	}
	updateTime := func(field string, value *time.Time) error {
		newValueString := r.FormValue(field)
		if newValueString == "" {
			return nil
		}
		newValue, err := time.Parse(pubdateLayout, newValueString)
		if err != nil {
			return err
		}
		*value = newValue
		return nil
	}

	updateString("ISBN", &book.ISBN)
	updateString("Title", &book.Title)
	updateString("Author", &book.Author)
	updateString("Description", &book.Description)
	updateString("CoverURL", &book.CoverURL)
	updateString("Publisher", &book.Publisher)
	if err := updateTime("PubDate", book.PubDate); err != nil {
		h.handleError(w, errors.New("invalid pubdate value"), http.StatusBadRequest)
		return
	}

	err = model.UpdateBook(h.db, book)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	book, err = model.GetBookByID(h.db, bookID)
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	h.handleSuccess(w, book)
}
