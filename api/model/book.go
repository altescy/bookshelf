package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Book struct {
	ID          uint64     `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" sql:"index"`
	ISBN        uint64     `json:"isbn"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Description string     `json:"description"`
	CoverURL    string     `json:"cover_url"`
	Publisher   string     `json:"publisher"`
	PubDate     time.Time  `json:"pubdate" gorm:"type:date"`
	Files       []File     `json:"files"`
}

func AddBook(db *gorm.DB, book *Book) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return handleBookError(tx.Create(book).Error)
	})
}

func GetBookByID(db *gorm.DB, bookID uint64) (*Book, error) {
	book := Book{}
	if err := db.First(&book, bookID).Error; err != nil {
		return nil, handleBookError(err)
	}
	return &book, nil
}

func GetBooks(db *gorm.DB) (*[]Book, error) {
	books := []Book{}
	if err := db.Find(&books).Error; err != nil {
		return nil, handleBookError(err)
	}
	return &books, nil
}

func GetBooksWithCount(db *gorm.DB, count uint64) (*[]Book, error) {
	books := []Book{}
	if err := db.Limit(count).Find(&books).Error; err != nil {
		return nil, handleBookError(err)
	}
	return &books, nil
}

func GetBooksWithNext(db *gorm.DB, next uint64) (*[]Book, error) {
	books := []Book{}
	if err := db.Where("id >= ?", next).Find(&books).Error; err != nil {
		return nil, handleBookError(err)
	}
	return &books, nil
}

func GetBooksWithNextCount(db *gorm.DB, next, count uint64) (*[]Book, error) {
	books := []Book{}
	if err := db.Where("id >= ?", next).Limit(count).Find(&books).Error; err != nil {
		return nil, handleBookError(err)
	}
	return &books, nil
}

func UpdateBook(db *gorm.DB, book *Book) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := db.Save(book).Error; err != nil {
			return handleBookError(err)
		}
		return nil
	})
}

func handleBookError(err error) error {
	if pgError, ok := err.(*pq.Error); ok {
		switch pgError.Code {
		case "23505":
			return ErrBookConflict
		}
	}

	switch {
	case gorm.IsRecordNotFoundError(err):
		return ErrBookNotFound
	default:
		return err
	}
}
