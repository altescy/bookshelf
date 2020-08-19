package storage

import (
	"io"
	"os"
	"path/filepath"
)

type FileSystemStorage struct {
	root string
	perm os.FileMode
}

func NewFileSystemStorage(root string, perm os.FileMode) *FileSystemStorage {
	return &FileSystemStorage{root: root, perm: perm}
}

func (s *FileSystemStorage) Upload(path string, body io.ReadSeeker) (err error) {
	path = filepath.Join(s.root, path)

	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, s.perm)
	if err != nil {
		return
	}

	f, err := os.Create(path)
	if err != nil {
		return
	}

	defer f.Close()

	_, err = io.Copy(f, body)
	return
}

func (s *FileSystemStorage) Download(w io.Writer, path string) (err error) {
	path = filepath.Join(s.root, path)

	f, err := os.Open(path)
	if err != nil {
		return
	}

	defer f.Close()

	_, err = io.Copy(w, f)
	return
}
