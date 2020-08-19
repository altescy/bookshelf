package model

import (
	"github.com/altescy/bookshelf/opds"
)

func EntriesFromBooks(books *[]Book) []opds.Entry {
	entries := make([]opds.Entry, 0, len(*books))

	for _, book := range *books {
		author := opds.Author{Name: book.Author}
		summary := opds.Summary{Type: "text", Text: book.Description}
		coverType, _ := MimeByFilename(book.CoverURL)
		links := []opds.Link{
			{Href: book.CoverURL, Type: coverType, Rel: opds.CoverRel},
		}
		for _, file := range book.Files {
			link := opds.Link{
				Href: file.Link,
				Type: file.MimeType,
				Rel:  opds.FileRel,
			}
			links = append(links, link)
		}
		entry := opds.Entry{
			ID:      "urn:uuid:" + book.UUID,
			Updated: book.UpdatedAt.UTC().Format(opds.AtomTime),
			Title:   book.Title,
			Author:  author,
			Summary: summary,
			Link:    links,
		}

		entries = append(entries, entry)
	}
	return entries
}
