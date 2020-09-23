package controller

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/altescy/bookshelf/model"
	"github.com/julienschmidt/httprouter"
)

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ext := "." + ps.ByName("ext")
	bookidString := ps.ByName("bookid")
	bookID, err := strconv.ParseUint(bookidString, 10, 64)
	if err != nil {
		h.handleError(w, errors.New("invalid bookid"), http.StatusBadRequest)
		return
	}

	mime, err := model.MimeByExt(ext)
	switch {
	case err == model.ErrMimeNotFound:
		h.handleError(w, err, http.StatusNotFound)
		return
	case err != nil:
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	err = model.DeleteFile(h.db, bookID, mime)
	switch {
	case err == model.ErrFileNotFound:
		h.handleError(w, err, http.StatusNotFound)
		return
	case err != nil:
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	h.handleSuccess(w, "successfully deleted")
}

func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ext := "." + ps.ByName("ext")
	bookidString := ps.ByName("bookid")
	bookID, err := strconv.ParseUint(bookidString, 10, 64)
	if err != nil {
		h.handleError(w, errors.New("invalid bookid"), http.StatusBadRequest)
		return
	}

	mime, err := model.MimeByExt(ext)
	switch {
	case err == model.ErrMimeNotFound:
		h.handleError(w, err, http.StatusNotFound)
		return
	case err != nil:
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	file, err := model.GetFile(h.db, bookID, mime)
	switch {
	case err == model.ErrFileNotFound:
		h.handleError(w, err, http.StatusNotFound)
		return
	case err != nil:
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", file.MimeType)
	w.WriteHeader(http.StatusOK)
	if err := h.storage.Download(w, file.Path); err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UploadFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookidString := ps.ByName("bookid")
	bookID, err := strconv.ParseUint(bookidString, 10, 64)
	if err != nil {
		h.handleError(w, errors.New("invalid bookid"), http.StatusBadRequest)
		return
	}

	_, err = model.GetBookByID(h.db, bookID)
	switch {
	case err == model.ErrBookNotFound:
		h.handleError(w, err, http.StatusNotFound)
		return
	case err != nil:
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		h.handleError(w, err, http.StatusInternalServerError)
		return
	}

	results := []map[string]interface{}{}

	for {
		file := model.File{BookID: bookID}

		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		filename := part.FileName()
		if filename == "" {
			continue
		}

		// set MIME type
		file.MimeType, err = model.MimeByFilename(filename)
		if err != nil {
			result := map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"content": err.Error(),
			}
			results = append(results, result)
			log.Printf("[ERROR] %+v", err)
			continue
		}

		// set file path
		mimeAlias, err := model.GetMimeAlias(file.MimeType)
		if err != nil {
			result := map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"content": err.Error(),
			}
			results = append(results, result)
			log.Printf("[ERROR] %+v", err)
			continue
		}
		file.Path = model.GenerateFilePath(bookID, mimeAlias)

		// read file
		b, err := ioutil.ReadAll(part)
		if err != nil {
			result := map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"content": err.Error(),
			}
			results = append(results, result)
			log.Printf("[ERROR] %+v", err)
			continue
		}

		// upload file to storage
		err = h.storage.Upload(file.Path, bytes.NewReader(b))
		if err != nil {
			result := map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"content": err.Error(),
			}
			results = append(results, result)
			log.Printf("[ERROR] %+v", err)
			continue
		}

		// add file to database
		err = model.AddFile(h.db, &file)
		if err != nil {
			result := map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"content": err.Error(),
			}
			results = append(results, result)
			log.Printf("[ERROR] %+v", err)
			continue
		}

		result := map[string]interface{}{
			"file":    filename,
			"status":  "ok",
			"content": file,
		}
		results = append(results, result)
	}

	h.handleSuccess(w, results)
}
