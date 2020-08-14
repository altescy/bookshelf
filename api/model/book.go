package model

import "github.com/jinzhu/gorm"

type Book struct {
	gorm.Model
	Title       string
	Author      string
	Description string
	CoverURL    string
	Files       []File
}

func AddBook(db *gorm.DB, book *Book) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(book).Error
	})
}
