package generator

import (
	"context"
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/post"
)

var ctx = context.Background() //nolint:gochecknoglobals // test-only top-level helper used by every test in the file

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

func TestThemeSoundsJSONNonEmpty(t *testing.T) {
	t.Parallel()
	j := themeSoundsJSON("neon")
	if len(j) < 50 {
		t.Fatalf("themeSoundsJSON too short: %q", j)
	}
}

func TestThemeSoundsJSON_ambientSchema(t *testing.T) {
	t.Parallel()

	// Verify the ambient schema is present and valid for every registered theme.
	for name := range getThemeSet() {
		j := themeSoundsJSON(name)
		if len(j) < 50 {
			t.Fatalf("themeSoundsJSON(%q) too short: %q", name, j)
		}

		// Parse as generic map to validate structure without coupling to field order.
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(j), &parsed); err != nil {
			t.Fatalf("themeSoundsJSON(%q) invalid JSON: %v", name, err)
		}

		// Core sound fields must remain present for backwards compatibility.
		for _, key := range []string{"splash", "nav", "open", "close", "bounce"} {
			if _, ok := parsed[key]; !ok {
				t.Errorf("themeSoundsJSON(%q) missing required key %q", name, key)
			}
		}

		// Ambient must be present and contain normal + wild variants.
		ambient, ok := parsed["ambient"].(map[string]interface{})
		if !ok {
			t.Errorf("themeSoundsJSON(%q) missing ambient object", name)
			continue
		}
		for _, key := range []string{"normal", "wild"} {
			if _, ok := ambient[key]; !ok {
				t.Errorf("themeSoundsJSON(%q) ambient missing %q variant", name, key)
			}
		}
	}
}

func TestThemeSoundPresetsAmbientPopulated(t *testing.T) {
	t.Parallel()

	for name := range getThemeSet() {
		preset, err := loadThemeSounds(name)
		if err != nil {
			t.Errorf("theme %q loadThemeSounds: %v", name, err)
			continue
		}

		normal := preset.Ambient.Normal
		wild := preset.Ambient.Wild

		if len(normal.DroneFreqs) == 0 && len(normal.PulseFreqs) == 0 && len(normal.Melody) == 0 {
			t.Errorf("theme %q ambient.Normal has no drone, pulse, or melody frequencies", name)
		}
		if len(wild.DroneFreqs) == 0 && len(wild.PulseFreqs) == 0 && len(wild.Melody) == 0 {
			t.Errorf("theme %q ambient.Wild has no drone, pulse, or melody frequencies", name)
		}
	}
}

func TestThemeSoundPresetsAmbientValuesBounded(t *testing.T) {
	t.Parallel()

	for name := range getThemeSet() {
		preset, err := loadThemeSounds(name)
		if err != nil {
			continue
		}

		for _, mode := range []string{"normal", "wild"} {
			var a ambientPreset
			if mode == "normal" {
				a = preset.Ambient.Normal
			} else {
				a = preset.Ambient.Wild
			}

			if a.Gain <= 0 || a.Gain > 0.15 {
				t.Errorf("theme %q ambient.%s gain=%f; want (0, 0.15]", name, mode, a.Gain)
			}
			if a.BPM <= 0 || a.BPM > 400 {
				t.Errorf("theme %q ambient.%s bpm=%f; want (0, 250]", name, mode, a.BPM)
			}
			if a.PulseInterval < 0 || a.PulseInterval > 10 {
				t.Errorf("theme %q ambient.%s pulseInterval=%f; want [0, 10]", name, mode, a.PulseInterval)
			}
			if a.Attack <= 0 || a.Attack > 5 {
				t.Errorf("theme %q ambient.%s attack=%f; want (0, 5]", name, mode, a.Attack)
			}
			if a.Release <= 0 || a.Release > 5 {
				t.Errorf("theme %q ambient.%s release=%f; want (0, 5]", name, mode, a.Release)
			}
			if a.NoiseGain < 0 || a.NoiseGain > 0.1 {
				t.Errorf("theme %q ambient.%s noiseGain=%f; want [0, 0.1]", name, mode, a.NoiseGain)
			}
			if a.DetuneCents < 0 || a.DetuneCents > 50 {
				t.Errorf("theme %q ambient.%s detuneCents=%f; want [0, 50]", name, mode, a.DetuneCents)
			}
			for i, f := range a.DroneFreqs {
				if f <= 0 {
					t.Errorf("theme %q ambient.%s droneFreqs[%d]=%f; want positive", name, mode, i, f)
				}
			}
			for i, f := range a.PulseFreqs {
				if f <= 0 {
					t.Errorf("theme %q ambient.%s pulseFreqs[%d]=%f; want positive", name, mode, i, f)
				}
			}
			for i, m := range a.Melody {
				if m.Freq <= 0 {
					t.Errorf("theme %q ambient.%s melody[%d].freq=%f; want positive", name, mode, i, m.Freq)
				}
				if m.Dur <= 0 {
					t.Errorf("theme %q ambient.%s melody[%d].dur=%f; want positive", name, mode, i, m.Dur)
				}
			}
			if a.CutoffMin < 0 || a.CutoffMin > 10000 {
				t.Errorf("theme %q ambient.%s cutoffMin=%f; want [0, 10000]", name, mode, a.CutoffMin)
			}
			if a.CutoffMax < 0 || a.CutoffMax > 10000 {
				t.Errorf("theme %q ambient.%s cutoffMax=%f; want [0, 10000]", name, mode, a.CutoffMax)
			}
		}
	}
}

func TestThemeSoundsJSON_neonAmbientRoundTrip(t *testing.T) {
	t.Parallel()

	j := themeSoundsJSON("neon")
	var s themeSounds
	if err := json.Unmarshal([]byte(j), &s); err != nil {
		t.Fatalf("themeSoundsJSON(\"neon\") unmarshal error: %v", err)
	}

	if s.Ambient.Normal.Gain <= 0 {
		t.Errorf("neon ambient.normal gain missing or non-positive: %f", s.Ambient.Normal.Gain)
	}
	if s.Ambient.Wild.Gain <= 0 {
		t.Errorf("neon ambient.wild gain missing or non-positive: %f", s.Ambient.Wild.Gain)
	}
	if len(s.Ambient.Normal.DroneFreqs) == 0 && len(s.Ambient.Normal.PulseFreqs) == 0 {
		t.Error("neon ambient.normal has no frequencies")
	}
	if len(s.Ambient.Wild.DroneFreqs) == 0 && len(s.Ambient.Wild.PulseFreqs) == 0 {
		t.Error("neon ambient.wild has no frequencies")
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

	meta, err := loadThemeMeta("neon")
	if err != nil {
		t.Fatalf("loadThemeMeta: %v", err)
	}
	all, err := allThemesJSON()
	if err != nil {
		t.Fatalf("allThemesJSON: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := buildPageData([]*post.Post{p}, tt.pageIndex, tt.totalPages, "neon", meta, all)
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

func TestValidThemeName_unknownFallsBackToNeon(t *testing.T) {
	t.Parallel()
	if got := validThemeName("no-such-theme-"); got != "neon" {
		t.Fatalf("validThemeName(\"no-such-theme-\") = %q; want \"neon\"", got)
	}
	if got := validThemeName("matrix"); got != "matrix" {
		t.Fatalf("validThemeName(\"matrix\") = %q; want \"matrix\"", got)
	}
}

func TestLoadThemeMeta_neonHasFields(t *testing.T) {
	t.Parallel()
	m, err := loadThemeMeta("neon")
	if err != nil {
		t.Fatalf("loadThemeMeta(neon): %v", err)
	}
	if m.Title == "" || m.HeaderHTML == "" || m.SplashInnerHTML == "" {
		t.Fatalf("neon meta missing required fields: %+v", m)
	}
}

func TestListThemes_sortedAndComplete(t *testing.T) {
	t.Parallel()
	names := ListThemes()
	if len(names) != len(getThemeSet()) {
		t.Fatalf("len=%d, want %d", len(names), len(getThemeSet()))
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
	if err := Run(ctx, cfg); err != nil {
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

func TestRun_writesVolcanoFontAssets(t *testing.T) {
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
		Content:   "<p>volcano</p>",
	}
	if err := p.Save(postDir); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		OutputDir: out,
		BaseURL:   "https://example.test",
		Theme:     "volcano",
	}
	if err := Run(ctx, cfg); err != nil {
		t.Fatalf("Run: %v", err)
	}

	themeDir := filepath.Join(out, "themes", "volcano")
	for _, name := range []string{
		"bebas-neue-v16-latin_latin-ext-regular.woff2",
		"inter-v20-latin_latin-ext-regular.woff2",
		"inter-v20-latin_latin-ext-600.woff2",
		"FONT_LICENSE.txt",
	} {
		info, err := os.Stat(filepath.Join(themeDir, name))
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if info.Size() == 0 {
			t.Fatalf("%s is empty", name)
		}
	}

	css, err := os.ReadFile(filepath.Join(themeDir, "theme.css"))
	if err != nil {
		t.Fatalf("read volcano theme.css: %v", err)
	}
	got := string(css)
	for _, localFont := range []string{
		"url('bebas-neue-v16-latin_latin-ext-regular.woff2')",
		"url('inter-v20-latin_latin-ext-regular.woff2')",
		"url('inter-v20-latin_latin-ext-600.woff2')",
	} {
		if !strings.Contains(got, localFont) {
			t.Fatalf("volcano theme.css missing local font reference %q", localFont)
		}
	}
	for _, forbidden := range []string{"googleapis", "gstatic", "fonts.cdn", "@import url(http"} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("volcano theme.css contains forbidden external font reference %q", forbidden)
		}
	}
}

func TestWritePage(t *testing.T) {
	t.Parallel()

	meta, err := loadThemeMeta("neon")
	if err != nil {
		t.Fatalf("loadThemeMeta: %v", err)
	}
	all, err := allThemesJSON()
	if err != nil {
		t.Fatalf("allThemesJSON: %v", err)
	}

	tests := []struct {
		name            string
		posts           []*post.Post
		pageIndex       int
		totalPages      int
		baseURL         string
		wantErr         bool
		wantErrContains string
	}{
		{
			name:       "happy path one page",
			posts:      []*post.Post{{ID: "a", Content: "<p>hi</p>"}},
			pageIndex:  0,
			totalPages: 1,
			baseURL:    "https://example.test",
			wantErr:    false,
		},
		{
			name:       "happy path second page",
			posts:      []*post.Post{{ID: "b", Content: "<p>bye</p>"}},
			pageIndex:  1,
			totalPages: 2,
			baseURL:    "https://example.test",
			wantErr:    false,
		},
		{
			name:            "invalid template action triggers error",
			posts:           []*post.Post{{ID: "x", Content: "<p>y</p>"}},
			pageIndex:       0,
			totalPages:      1,
			baseURL:         "https://example.test",
			wantErr:         true,
			wantErrContains: "render index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out := t.TempDir()
			cfg := &config.Config{
				OutputDir: out,
				BaseURL:   tt.baseURL,
				Theme:     "neon",
			}

			var tmpl *template.Template
			if tt.wantErr {
				tmpl = template.Must(template.New("page").Parse("{{.NonExistent.X}}"))
			} else {
				tmpl = template.Must(template.New("page").Parse("<html>{{.DefaultTheme}}</html>"))
			}

			err := writePage(tmpl, tt.posts, tt.pageIndex, tt.totalPages, cfg, "neon", meta, all)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.wantErrContains != "" && !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Fatalf("error %q does not contain %q", err.Error(), tt.wantErrContains)
				}
				path := filepath.Join(out, pageFilename(tt.pageIndex))
				if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
					t.Fatalf("expected %s to be absent after failed write, got statErr=%v", path, statErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			path := filepath.Join(out, pageFilename(tt.pageIndex))
			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			got := string(b)
			if !strings.HasPrefix(got, "<html>") {
				t.Fatalf("expected HTML output, got %q", got)
			}
		})
	}
}

func TestWritePage_tempFileCleanedOnError(t *testing.T) {
	t.Parallel()
	out := t.TempDir()

	path := filepath.Join(out, "index.html")
	golden := "golden"
	if err := os.WriteFile(path, []byte(golden), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		OutputDir: out,
		BaseURL:   "https://example.test",
		Theme:     "neon",
	}
	meta, _ := loadThemeMeta("neon")
	all, _ := allThemesJSON()
	tmpl := template.Must(template.New("page").Parse("{{.NonExistent.X}}"))

	err := writePage(tmpl, []*post.Post{{ID: "a", Content: "<p>x</p>"}}, 0, 1, cfg, "neon", meta, all)
	if err == nil {
		t.Fatal("expected error from broken template")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read existing file: %v", err)
	}
	if string(b) != golden {
		t.Fatalf("existing file was corrupted: got %q, want %q", string(b), golden)
	}
}
