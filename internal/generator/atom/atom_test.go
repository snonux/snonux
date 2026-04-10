package atom

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/post"
)

func TestGenerate_writesAtomXML(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := &config.Config{
		OutputDir: dir,
		BaseURL:   "https://example.test",
	}
	posts := []*post.Post{
		{
			ID:        "p1",
			Timestamp: time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC),
			Content:   "<p>hello</p>",
		},
	}

	if err := Generate(posts, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "atom.xml"))
	if err != nil {
		t.Fatalf("read atom.xml: %v", err)
	}
	s := string(data)
	if !strings.Contains(s, `xmlns="http://www.w3.org/2005/Atom"`) {
		t.Fatalf("missing atom xmlns: %s", s)
	}
	if !strings.Contains(s, "https://example.test/posts/p1/") {
		t.Fatalf("missing entry link: %s", s)
	}
	if !strings.Contains(s, "hello") || !strings.Contains(s, `type="html"`) {
		t.Fatalf("missing content: %s", s)
	}
}

func TestGenerate_emptyPosts(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := &config.Config{OutputDir: dir, BaseURL: "https://x.test"}

	if err := Generate(nil, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "atom.xml"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var feed struct {
		Entries []struct {
			Title string `xml:"title"`
		} `xml:"entry"`
	}
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("xml: %v", err)
	}
	if len(feed.Entries) != 0 {
		t.Fatalf("want 0 entries, got %d", len(feed.Entries))
	}
}

func TestGenerate_limitPostsPerPage(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfg := &config.Config{OutputDir: dir, BaseURL: "https://x.test"}

	var posts []*post.Post
	for i := 0; i < 50; i++ {
		posts = append(posts, &post.Post{
			ID:        "id",
			Timestamp: time.Date(2026, 1, i+1, 12, 0, 0, 0, time.UTC),
			Content:   "c",
		})
	}

	if err := Generate(posts, cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(dir, "atom.xml"))
	var feed struct {
		Entries []struct{} `xml:"entry"`
	}
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("xml: %v", err)
	}
	if len(feed.Entries) != config.PostsPerPage {
		t.Fatalf("want %d entries, got %d", config.PostsPerPage, len(feed.Entries))
	}
}
