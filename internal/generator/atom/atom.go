// Package atom serializes the site Atom 1.0 feed (atom.xml). It depends only on
// internal/config and internal/post. It does not import HTML themes, shared nav
// templates, or the html/template page pipeline — those live in the parent
// generator package.
package atom

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/post"
)

// feed is the root element of an Atom 1.0 feed document.
type feed struct {
	XMLName xml.Name `xml:"feed"`
	XMLNS   string   `xml:"xmlns,attr"`
	Title   string   `xml:"title"`
	Link    link     `xml:"link"`
	Updated string   `xml:"updated"`
	ID      string   `xml:"id"`
	Entries []entry  `xml:"entry"`
}

type link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

type entry struct {
	Title   string  `xml:"title"`
	Link    link    `xml:"link"`
	ID      string  `xml:"id"`
	Updated string  `xml:"updated"`
	Content content `xml:"content"`
}

type content struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// Generate writes atom.xml to cfg.OutputDir containing the most recent
// min(len(posts), config.PostsPerPage) entries. Posts must already be sorted
// newest-first (as produced by generator.Run).
func Generate(posts []*post.Post, cfg *config.Config) error {
	limit := config.PostsPerPage
	if len(posts) < limit {
		limit = len(posts)
	}

	recent := posts[:limit]
	entries := buildEntries(recent, cfg.BaseURL)

	updated := time.Now().UTC().Format(time.RFC3339)
	if len(recent) > 0 {
		updated = recent[0].Timestamp.UTC().Format(time.RFC3339)
	}

	f := feed{
		XMLNS:   "http://www.w3.org/2005/Atom",
		Title:   "snonux.foo",
		Link:    link{Href: cfg.BaseURL + "/"},
		Updated: updated,
		ID:      cfg.BaseURL + "/",
		Entries: entries,
	}

	return writeFile(f, filepath.Join(cfg.OutputDir, "atom.xml"))
}

func buildEntries(posts []*post.Post, baseURL string) []entry {
	entries := make([]entry, 0, len(posts))

	for _, p := range posts {
		entryURL := fmt.Sprintf("%s/posts/%s/", baseURL, p.ID)
		entries = append(entries, entry{
			Title:   fmt.Sprintf("Post %s", p.ID),
			Link:    link{Href: entryURL, Rel: "alternate"},
			ID:      entryURL,
			Updated: p.Timestamp.UTC().Format(time.RFC3339),
			Content: content{Type: "html", Value: p.Content},
		})
	}

	return entries
}

func writeFile(f feed, path string) error {
	data, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal atom feed: %w", err)
	}

	content := append([]byte(xml.Header), data...)

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write atom.xml: %w", err)
	}

	return nil
}
