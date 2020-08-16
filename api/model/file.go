package model

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/altescy/bookshelf/api/storage"
	"github.com/jinzhu/gorm"
	"github.com/oklog/ulid"
)

type File struct {
	ID        uint64     `json:"ID"`
	CreatedAt time.Time  `json:"CreatedAt"`
	UpdatedAt time.Time  `json:"UpdatedAt"`
	DeletedAt *time.Time `json:"-" sql:"index"`
	BookID    uint64     `json:"BookId"`
	MimeType  string     `json:"MimeType"`
	Path      string     `json:"Path"`
}

func AddFile(db *gorm.DB, storage storage.Storage, file *File, body io.ReadSeeker) error {
	// check book existence
	if _, err := GetBookByID(db, file.BookID); err != nil {
		return err
	}

	// check the same MimeType existence
	err := db.Take(&File{}, "book_id=? and mime_type=?", file.BookID, file.MimeType).Error
	switch {
	case err == nil:
		return ErrFileConflict
	case !gorm.IsRecordNotFoundError(err):
		return err
	}

	// set file path
	filename := getULID()
	mimealias, err := GetMimeAlias(file.MimeType)
	if err != nil {
		return err
	}
	file.Path = fmt.Sprintf("%d/%s/%s", file.BookID, mimealias, filename)

	// upload file to storage
	if err := storage.Upload(file.Path, body); err != nil {
		return err
	}

	// add file to database
	return db.Transaction(func(tx *gorm.DB) error {
		return handleBookError(tx.Save(file).Error)
	})
}

func getULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
