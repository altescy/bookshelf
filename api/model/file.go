package model

type File struct {
	ID       uint64 `json:"ID"`
	BookID   uint64 `json:"BookId" gorm:"primary_key"`
	MimeType string `json:"MimeType" gorm:"primary_key"`
	Path     string `json:"Path"`
}
