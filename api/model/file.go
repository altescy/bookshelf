package model

import (
	"fmt"
	"math/rand"
	"time"

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

func AddFile(db *gorm.DB, file *File) error {
	// check the same MimeType existence
	err := db.Take(&File{}, "book_id=? and mime_type=?", file.BookID, file.MimeType).Error
	switch {
	case err == nil:
		return ErrFileConflict
	case !gorm.IsRecordNotFoundError(err):
		return err
	}

	// add file to database
	return db.Transaction(func(tx *gorm.DB) error {
		return handleBookError(tx.Save(file).Error)
	})
}

func GetFile(db *gorm.DB, bookID uint64, mime string) (*File, error) {
	file := File{}
	err := db.Last(&file, "book_id=? and mime_type=?", bookID, mime).Error
	switch {
	case gorm.IsRecordNotFoundError(err):
		return nil, ErrFileNotFound
	case err != nil:
		return nil, err
	}
	return &file, nil
}

func GenerateFilePath(bookID uint64, mimeAlias string) string {
	filename := generateULID()
	path := fmt.Sprintf("%d/%s/%s", bookID, mimeAlias, filename)
	return path
}

func generateULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
