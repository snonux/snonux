package generator

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/post"
)

func TestPageFilename(t *testing.T) {
	t.Parallel()

	tests := []struct {
		index int
		want  string
	}{
		{0, "index.html"},
		{1, "page2.html"},
		{2, "page3.html"},
	}

	for _, tt := range tests {
		if got := pageFilename(tt.index); got != tt.want {
			t.Fatalf("pageFilename(%d) = %q; want %q", tt.index, got, tt.want)
		}
	}
}

func TestPaginate(t *testing.T) {
	t.Parallel()

	p := func(ids ...string) []*post.Post {
		out := make([]*post.Post, len(ids))
		for i, id := range ids {
			out[i] = &post.Post{ID: id}
		}
		return out
	}

	tests := []struct {
		name     string
		posts    []*post.Post
		pageSize int
		wantLens []int
	}{
		{name: "empty", posts: nil, pageSize: 3, wantLens: nil},
		{name: "one page exact", posts: p("a", "b"), pageSize: 2, wantLens: []int{2}},
		{name: "two pages", posts: p("a", "b", "c"), pageSize: 2, wantLens: []int{2, 1}},
		{name: "singleton pages", posts: p("x", "y"), pageSize: 1, wantLens: []int{1, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pages := paginate(tt.posts, tt.pageSize)
			if len(pages) != len(tt.wantLens) {
				t.Fatalf("len(pages)=%d; want %d", len(pages), len(tt.wantLens))
			}
			for i, n := range tt.wantLens {
				if len(pages[i]) != n {
					t.Fatalf("page %d len=%d; want %d", i, len(pages[i]), n)
				}
			}
		})
	}
}

func TestJSONStringOrNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want template.JS
	}{
		{in: "", want: "null"},
		{in: "page2.html", want: `"page2.html"`},
		{in: `say "hi"`, want: `"say \"hi\""`},
	}

	for _, tt := range tests {
		got := jsonStringOrNull(tt.in)
		if got != tt.want {
			t.Fatalf("jsonStringOrNull(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

func TestThemeSoundPresetsMatchRegistry(t *testing.T) {
	t.Parallel()
	for name := range themeRegistry {
		if _, ok := themeSoundPresets[name]; !ok {
			t.Errorf("theme %q has no sound preset in themeSoundPresets", name)
		}
	}
}

func TestThemeSoundsJSONNonEmpty(t *testing.T) {
	t.Parallel()
	j := themeSoundsJSON("neon")
	if len(j) < 50 {
		t.Fatalf("themeSoundsJSON too short: %q", j)
	}
}

func TestFormatPostTime(t *testing.T) {
	t.Parallel()

	tm := time.Date(2026, 4, 9, 14, 30, 0, 0, time.FixedZone("CET", 3600))
	got := formatPostTime(tm)
	want := "09.04.26 • 13:30 UTC"
	if got != want {
		t.Fatalf("formatPostTime = %q; want %q", got, want)
	}
}

func TestBuildPageData_navLinks(t *testing.T) {
	t.Parallel()

	p := &post.Post{
		ID:        "1",
		Timestamp: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
		Content:   "<p>x</p>",
	}

	tests := []struct {
		name           string
		pageIndex      int
		totalPages     int
		wantPrev       string
		wantNext       string
		wantPrevJSON   template.JS
		wantNextJSON   template.JS
		wantPostsCount int
	}{
		{
			name:           "first of three",
			pageIndex:      0,
			totalPages:     3,
			wantPrev:       "",
			wantNext:       "page2.html",
			wantPrevJSON:   "null",
			wantNextJSON:   `"page2.html"`,
			wantPostsCount: 1,
		},
		{
			name:           "middle",
			pageIndex:      1,
			totalPages:     3,
			wantPrev:       "index.html?splash=0",
			wantNext:       "page3.html",
			wantPrevJSON:   `"index.html?splash=0"`,
			wantNextJSON:   `"page3.html"`,
			wantPostsCount: 1,
		},
		{
			name:           "last",
			pageIndex:      2,
			totalPages:     3,
			wantPrev:       "page2.html",
			wantNext:       "",
			wantPrevJSON:   `"page2.html"`,
			wantNextJSON:   "null",
			wantPostsCount: 1,
		},
		{
			name:           "single page",
			pageIndex:      0,
			totalPages:     1,
			wantPrev:       "",
			wantNext:       "",
			wantPrevJSON:   "null",
			wantNextJSON:   "null",
			wantPostsCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := buildPageData([]*post.Post{p}, tt.pageIndex, tt.totalPages, "neon")
			if data.PrevPage != tt.wantPrev {
				t.Fatalf("PrevPage=%q; want %q", data.PrevPage, tt.wantPrev)
			}
			if data.NextPage != tt.wantNext {
				t.Fatalf("NextPage=%q; want %q", data.NextPage, tt.wantNext)
			}
			if data.PrevPageJSON != tt.wantPrevJSON {
				t.Fatalf("PrevPageJSON=%q; want %q", data.PrevPageJSON, tt.wantPrevJSON)
			}
			if data.NextPageJSON != tt.wantNextJSON {
				t.Fatalf("NextPageJSON=%q; want %q", data.NextPageJSON, tt.wantNextJSON)
			}
			if len(data.Posts) != tt.wantPostsCount {
				t.Fatalf("len(Posts)=%d", len(data.Posts))
			}
		})
	}
}

func TestGetTheme_unknownFallsBackToNeon(t *testing.T) {
	t.Parallel()
	if got, want := getTheme("no-such-theme-"), getTheme("neon"); got != want {
		t.Fatal("expected neon fallback")
	}
}

func TestInjectSharedHead_addsFaviconLink(t *testing.T) {
	t.Parallel()

	got := injectSharedHead(getTheme("neon"))
	if !strings.Contains(got, `rel="icon" href="favicon.ico"`) {
		t.Fatalf("favicon link missing from theme head")
	}
}

func TestListThemes_sortedAndComplete(t *testing.T) {
	t.Parallel()
	names := ListThemes()
	if len(names) != len(themeRegistry) {
		t.Fatalf("len=%d, want %d", len(names), len(themeRegistry))
	}
	for i := 1; i < len(names); i++ {
		if names[i] <= names[i-1] {
			t.Fatalf("not strictly sorted: %v", names)
		}
	}
}

func TestLoadAllPosts_missingPostsDir(t *testing.T) {
	t.Parallel()
	posts, err := loadAllPosts(t.TempDir())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if posts != nil {
		t.Fatalf("want nil slice, got %v", posts)
	}
}

func TestRun_writesPagesAndAtom(t *testing.T) {
	t.Parallel()

	out := t.TempDir()
	postDir := filepath.Join(out, "posts", "a1")
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		t.Fatal(err)
	}
	p := &post.Post{
		ID:        "a1",
		Timestamp: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
		PostType:  post.TypeText,
		Content:   "<p>hello</p>",
	}
	if err := p.Save(postDir); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		OutputDir: out,
		BaseURL:   "https://example.test",
		Theme:     "neon",
	}
	if err := Run(cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "index.html")); err != nil {
		t.Fatalf("index.html: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "atom.xml")); err != nil {
		t.Fatalf("atom.xml: %v", err)
	}
	if _, err := os.Stat(filepath.Join(out, "favicon.ico")); err != nil {
		t.Fatalf("favicon.ico: %v", err)
	}
	indexHTML, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatalf("read index.html: %v", err)
	}
	if !strings.Contains(string(indexHTML), `rel="icon" href="favicon.ico"`) {
		t.Fatalf("index.html missing favicon link: %s", string(indexHTML))
	}
}
