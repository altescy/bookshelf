package opds

import (
	"encoding/xml"
	"time"
)

const (
	AtomTime = "2006-01-02T15:04:05Z"
	DirMime  = "application/atom+xml;profile=opds-catalog;kind=navigation"
	DirRel   = "subsection"
	FileRel  = "http://opds-spec.org/acquisition"
	CoverRel = "http://opds-spec.org/cover"
)

// Feed is a main frame of OPDS.
type Feed struct {
	XMLName xml.Name `xml:"feed"`
	ID      string   `xml:"id"`
	Title   string   `xml:"title"`
	Xmlns   string   `xml:"xmlns,attr"`
	Updated string   `xml:"updated"`
	Link    []Link   `xml:"link"`
	Entry   []Entry  `xml:"entry"`
}

// Link is link properties.
type Link struct {
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
	Rel  string `xml:"rel,attr,ommitempty"`
}

// Entry is a struct of OPDS entry properties.
type Entry struct {
	ID      string  `xml:"id"`
	Updated string  `xml:"updated"`
	Title   string  `xml:"title"`
	Author  Author  `xml:"author"`
	Summary Summary `xml:"summary"`
	Link    []Link  `xml:"link"`
}

type Author struct {
	Name string `xml:"name"`
}

type Summary struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

func BuildFeed(id, title, href string, entries []Entry) *Feed {
	return &Feed{
		ID:      id,
		Title:   title,
		Xmlns:   "http://www.w3.org/2005/Atom",
		Updated: time.Now().UTC().Format(AtomTime),
		Link: []Link{
			{
				Href: href,
				Type: DirMime,
				Rel:  "start",
			},
			{
				Href: href,
				Type: DirMime,
				Rel:  "self",
			},
		},
		Entry: entries,
	}
}
