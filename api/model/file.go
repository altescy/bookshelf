package model

import (
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	BookID   uint
	Path     string `gorm:"unique_index"`
	MimeType string
}

// func AddFile
//     - upload file to s3
//         - file path is /bucket/dir/BookID/MimeType
//     - register file into DB
