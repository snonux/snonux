package generator

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/post"
)

// atomFeed is the root element of an Atom 1.0 feed document.
type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	XMLNS   string      `xml:"xmlns,attr"`
	Title   string      `xml:"title"`
	Link    atomLink    `xml:"link"`
	Updated string      `xml:"updated"`
	ID      string      `xml:"id"`
	Entries []atomEntry `xml:"entry"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

type atomEntry struct {
	Title   string   `xml:"title"`
	Link    atomLink `xml:"link"`
	ID      string   `xml:"id"`
	Updated string   `xml:"updated"`
	Content atomContent `xml:"content"`
}

type atomContent struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// generateAtom writes atom.xml to cfg.OutputDir containing the most recent
// min(len(posts), config.PostsPerPage) entries.
func generateAtom(posts []*post.Post, cfg *config.Config) error {
	limit := config.PostsPerPage
	if len(posts) < limit {
		limit = len(posts)
	}

	recent := posts[:limit]
	entries := buildAtomEntries(recent, cfg.BaseURL)

	updated := time.Now().UTC().Format(time.RFC3339)
	if len(recent) > 0 {
		updated = recent[0].Timestamp.UTC().Format(time.RFC3339)
	}

	feed := atomFeed{
		XMLNS:   "http://www.w3.org/2005/Atom",
		Title:   "snonux.foo",
		Link:    atomLink{Href: cfg.BaseURL + "/"},
		Updated: updated,
		ID:      cfg.BaseURL + "/",
		Entries: entries,
	}

	return writeAtomFile(feed, filepath.Join(cfg.OutputDir, "atom.xml"))
}

// buildAtomEntries converts a slice of posts into Atom entry elements.
func buildAtomEntries(posts []*post.Post, baseURL string) []atomEntry {
	entries := make([]atomEntry, 0, len(posts))

	for _, p := range posts {
		entryURL := fmt.Sprintf("%s/posts/%s/", baseURL, p.ID)
		entry := atomEntry{
			Title:   fmt.Sprintf("Post %s", p.ID),
			Link:    atomLink{Href: entryURL, Rel: "alternate"},
			ID:      entryURL,
			Updated: p.Timestamp.UTC().Format(time.RFC3339),
			Content: atomContent{Type: "html", Value: p.Content},
		}
		entries = append(entries, entry)
	}

	return entries
}

// writeAtomFile marshals feed to XML and writes it to path with XML declaration.
func writeAtomFile(feed atomFeed, path string) error {
	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal atom feed: %w", err)
	}

	content := append([]byte(xml.Header), data...)

	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write atom.xml: %w", err)
	}

	return nil
}
