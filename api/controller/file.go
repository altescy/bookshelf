package controller

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/altescy/bookshelf/api/model"
	"github.com/julienschmidt/httprouter"
)

func (h *Handler) UploadFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bookidString := ps.ByName("bookid")
	bookID, err := strconv.ParseUint(bookidString, 10, 64)
	if err != nil {
		h.handleError(w, errors.New("invalid bookid"), http.StatusBadRequest)
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

		file.MimeType, err = model.MimeByFilename(filename)
		if err != nil {
			result := map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"content": err.Error(),
			}
			results = append(results, result)
			continue
		}

		b, err := ioutil.ReadAll(part)
		if err != nil {
			result := map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"content": err.Error(),
			}
			results = append(results, result)
			continue
		}

		err = model.AddFile(h.db, h.storage, &file, bytes.NewReader(b))
		if err != nil {
			result := map[string]interface{}{
				"file":    filename,
				"status":  "error",
				"content": err.Error(),
			}
			results = append(results, result)
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
