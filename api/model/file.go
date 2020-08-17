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
	BookID    uint64     `json:"BookID"`
	MimeType  string     `json:"MimeType"`
	Path      string     `json:"-"`
	Link      string     `json:"Link" gorm:"-"`
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
		return handleFileError(tx.Save(&file).Error)
	})
}

func DeleteFile(db *gorm.DB, bookID uint64, mime string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := db.Delete(File{}, "book_id=? and mime_type=?", bookID, mime).Error
		return handleFileError(err)
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

func handleFileError(err error) error {
	switch {
	case gorm.IsRecordNotFoundError(err):
		return ErrFileNotFound
	default:
		return err
	}
}
