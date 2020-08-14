package model

type File struct {
	ID       uint64 `json:"id"`
	BookID   uint64 `json:"book_id" gorm:"primary_key"`
	MimeType string `json:"mimetype" gorm:"primary_key"`
	Path     string `json:"path"`
}
