package storage

import "io"

type Storage interface {
	Upload(path string, body io.ReadSeeker) error
	Download(w io.Writer, path string) error
}
