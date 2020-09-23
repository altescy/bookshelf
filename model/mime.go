package model

import (
	"mime"
	"path/filepath"
	"strings"
)

var extToMime = map[string]string{
	".azw3": "application/x-mobi8-ebook",
	".epub": "application/epub+zip",
	".fb2":  "application/fb2+zip",
	".mobi": "application/x-mobipocket-ebook",
	".pdf":  "application/pdf",
	".txt":  "text/plain",
}

var MimeAlias = map[string]string{
	"application/x-mobi8-ebook":      "azw3",
	"application/epub+zip":           "epub",
	"application/fb2+zip":            "fb2",
	"application/x-mobipocket-ebook": "mobi",
	"application/pdf":                "pdf",
	"text/plain":                     "txt",
}

func GetMimeAlias(mime string) (string, error) {
	alias := MimeAlias[mime]
	if alias == "" {
		return "", ErrMimeNotFound
	}
	return alias, nil

}

func GetMimeAliasByFilename(filename string) (string, error) {
	mime, err := MimeByFilename(filename)
	if err != nil {
		return "", nil
	}

	mimeAlias, err := GetMimeAlias(mime)
	if err != nil {
		return "", nil
	}

	return mimeAlias, nil
}

func GetMimes() map[string]string {
	return copyMimes(extToMime)
}

func MimeByExt(ext string) (string, error) {
	mime := extToMime[ext]
	if mime == "" {
		return "", ErrMimeNotFound
	}
	return mime, nil
}

func MimeByFilename(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	mimeType, ok := extToMime[ext]
	if ok {
		return mimeType, nil
	}

	mimeType = mime.TypeByExtension(ext)
	if mimeType == "" {
		return "", ErrInvalidExt
	}
	// The mimeType returned by TypeByExtension may contain
	// some arguments like "text/plain; charset=utf-8". But
	// in this statement, we ignore such arguments.
	mimeType = strings.Split(mimeType, ";")[0]

	return mimeType, nil
}

func copyMimes(m map[string]string) map[string]string {
	cp := map[string]string{}
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
