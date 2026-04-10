// Package integrationtests runs end-to-end tests of the snonux generator pipeline.
// Each test creates temporary input/output directories, places fixture files, runs
// the full processor+generator pipeline, and asserts the expected outputs.
package integrationtests

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/generator"
	"codeberg.org/snonux/snonux/internal/processor"
)

// runPipeline executes both pipeline stages and returns the config used.
func runPipeline(t *testing.T, inputDir, outputDir string) *config.Config {
	t.Helper()

	cfg := &config.Config{
		InputDir:  inputDir,
		OutputDir: outputDir,
		BaseURL:   "https://snonux.foo",
		Theme:     "neon",
	}

	_, err := processor.Run(cfg)
	if err != nil {
		t.Fatalf("processor.Run: %v", err)
	}

	if err := generator.Run(cfg); err != nil {
		t.Fatalf("generator.Run: %v", err)
	}

	return cfg
}

// makeDirs creates temporary input and output directories for a test.
func makeDirs(t *testing.T) (inputDir, outputDir string) {
	t.Helper()

	base := t.TempDir()
	inputDir = filepath.Join(base, "inbox")
	outputDir = filepath.Join(base, "outdir")

	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatal(err)
	}

	return inputDir, outputDir
}

// readFile is a helper that reads a file and fails the test on error.
func readFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	return string(data)
}

// assertContains fails the test if content does not contain substr.
func assertContains(t *testing.T, content, substr, label string) {
	t.Helper()

	if !strings.Contains(content, substr) {
		t.Errorf("%s: expected to contain %q\ngot:\n%s", label, substr, content[:min(len(content), 500)])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestTxtInput verifies plain text files are converted to posts.
func TestTxtInput(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	if err := os.WriteFile(filepath.Join(inputDir, "hello.txt"), []byte("Hello, Nexus!"), 0o644); err != nil {
		t.Fatal(err)
	}

	runPipeline(t, inputDir, outputDir)

	// Source file should have been removed after processing.
	if _, err := os.Stat(filepath.Join(inputDir, "hello.txt")); !os.IsNotExist(err) {
		t.Error("source file should have been deleted from input dir")
	}

	// A post directory should exist under outdir/posts/.
	entries, err := os.ReadDir(filepath.Join(outputDir, "posts"))
	if err != nil {
		t.Fatalf("read posts dir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 post dir, got %d", len(entries))
	}

	// index.html must contain the post text.
	index := readFile(t, filepath.Join(outputDir, "index.html"))
	assertContains(t, index, "Hello, Nexus!", "index.html")
}

// TestMarkdownInput verifies Markdown files are converted to HTML.
func TestMarkdownInput(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	md := "# Hello Nexus\n\nThis is **bold** text."
	if err := os.WriteFile(filepath.Join(inputDir, "post.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}

	runPipeline(t, inputDir, outputDir)

	index := readFile(t, filepath.Join(outputDir, "index.html"))
	assertContains(t, index, "<strong>bold</strong>", "index.html markdown bold")
	assertContains(t, index, "<h1>", "index.html markdown h1")
}

// assertStandaloneImagePost checks index.html and posts/<id>/image.jpg after a lone image input.
func assertStandaloneImagePost(t *testing.T, outputDir string) {
	t.Helper()

	index := readFile(t, filepath.Join(outputDir, "index.html"))
	assertContains(t, index, `<img`, "index.html image tag")
	assertContains(t, index, `image.jpg`, "index.html image filename")

	postDirs, err := os.ReadDir(filepath.Join(outputDir, "posts"))
	if err != nil {
		t.Fatalf("read posts dir: %v", err)
	}
	if len(postDirs) != 1 {
		t.Fatalf("expected 1 post, got %d", len(postDirs))
	}
	imgPath := filepath.Join(outputDir, "posts", postDirs[0].Name(), "image.jpg")
	if _, err := os.Stat(imgPath); err != nil {
		t.Errorf("expected image.jpg in post dir: %v", err)
	}
}

// TestPNGInput verifies .png files are converted to JPEG posts and embedded in pages.
func TestPNGInput(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	writeSamplePNG(t, filepath.Join(inputDir, "photo.png"))
	runPipeline(t, inputDir, outputDir)
	assertStandaloneImagePost(t, outputDir)
}

// TestJPGInput verifies .jpg files are processed the same way as PNG.
func TestJPGInput(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	writeSampleJPEG(t, filepath.Join(inputDir, "photo.jpg"))
	runPipeline(t, inputDir, outputDir)
	assertStandaloneImagePost(t, outputDir)
}

// TestJPEGInput verifies the .jpeg extension is accepted.
func TestJPEGInput(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	writeSampleJPEG(t, filepath.Join(inputDir, "snapshot.jpeg"))
	runPipeline(t, inputDir, outputDir)
	assertStandaloneImagePost(t, outputDir)
}

// TestGIFInput verifies .gif files are decoded (first frame) and output as JPEG.
func TestGIFInput(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	writeSampleGIF(t, filepath.Join(inputDir, "anim.gif"))
	runPipeline(t, inputDir, outputDir)
	assertStandaloneImagePost(t, outputDir)
}

// TestAudioInput verifies .mp3 files are copied and an audio element is generated.
func TestAudioInput(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	// Write a minimal non-empty file as a stand-in for MP3 content.
	if err := os.WriteFile(filepath.Join(inputDir, "track.mp3"), []byte("ID3fake"), 0o644); err != nil {
		t.Fatal(err)
	}

	runPipeline(t, inputDir, outputDir)

	index := readFile(t, filepath.Join(outputDir, "index.html"))
	assertContains(t, index, `<audio`, "index.html audio tag")
	assertContains(t, index, `track.mp3`, "index.html audio filename")
}

// TestMarkdownWithImage verifies that a Markdown post referencing a local image
// copies the image into the post dir and updates the src path.
func TestMarkdownWithImage(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	md := "Look at this:\n\n![cool pic](photo.png)\n"
	if err := os.WriteFile(filepath.Join(inputDir, "post.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}

	writeSamplePNG(t, filepath.Join(inputDir, "photo.png"))

	runPipeline(t, inputDir, outputDir)

	postDirs, _ := os.ReadDir(filepath.Join(outputDir, "posts"))
	if len(postDirs) != 1 {
		t.Fatalf("expected 1 post, got %d", len(postDirs))
	}

	// The referenced image should be copied into the post dir.
	imgPath := filepath.Join(outputDir, "posts", postDirs[0].Name(), "photo.png")
	if _, err := os.Stat(imgPath); err != nil {
		t.Errorf("expected photo.png in post dir: %v", err)
	}
}

// TestPagination verifies that 45 posts are split across two pages (42 + 3).
func TestPagination(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	for i := 0; i < 45; i++ {
		name := fmt.Sprintf("post%02d.txt", i)
		content := fmt.Sprintf("Post number %d", i)
		if err := os.WriteFile(filepath.Join(inputDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	runPipeline(t, inputDir, outputDir)

	// index.html should exist and contain 42 posts.
	index := readFile(t, filepath.Join(outputDir, "index.html"))
	if count := strings.Count(index, `class="post"`); count != 42 {
		t.Errorf("index.html: expected 42 posts, got %d", count)
	}

	// page2.html should exist and contain 3 posts.
	page2 := readFile(t, filepath.Join(outputDir, "page2.html"))
	if count := strings.Count(page2, `class="post"`); count != 3 {
		t.Errorf("page2.html: expected 3 posts, got %d", count)
	}
}

// TestPaginationNavLinks verifies prev/next navigation links are positioned correctly.
func TestPaginationNavLinks(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	for i := 0; i < 45; i++ {
		if err := os.WriteFile(filepath.Join(inputDir, fmt.Sprintf("p%02d.txt", i)), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	runPipeline(t, inputDir, outputDir)

	index := readFile(t, filepath.Join(outputDir, "index.html"))
	// index.html (page 1) has no prev, should have next link (page2.html).
	assertContains(t, index, "page2.html", "index.html next link")
	if strings.Contains(index, "NEWER TRANSMISSIONS") {
		t.Error("index.html should not have a prev-page link")
	}

	page2 := readFile(t, filepath.Join(outputDir, "page2.html"))
	// page2.html should have a prev link (index.html) and no next.
	assertContains(t, page2, "NEWER TRANSMISSIONS", "page2.html prev link")
	if strings.Contains(page2, "OLDER TRANSMISSIONS") {
		t.Error("page2.html should not have a next-page link")
	}
}

// TestAtomFeed verifies that atom.xml is well-formed and contains ≤42 entries.
func TestAtomFeed(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	for i := 0; i < 5; i++ {
		if err := os.WriteFile(filepath.Join(inputDir, fmt.Sprintf("p%d.txt", i)), []byte("feed post"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	runPipeline(t, inputDir, outputDir)

	atomPath := filepath.Join(outputDir, "atom.xml")
	data, err := os.ReadFile(atomPath)
	if err != nil {
		t.Fatalf("read atom.xml: %v", err)
	}

	// Validate well-formed XML.
	var feed struct {
		XMLName xml.Name `xml:"feed"`
		Entries []struct {
			Title string `xml:"title"`
		} `xml:"entry"`
	}
	if err := xml.Unmarshal(data, &feed); err != nil {
		t.Fatalf("atom.xml not valid XML: %v", err)
	}

	if len(feed.Entries) != 5 {
		t.Errorf("expected 5 entries in atom.xml, got %d", len(feed.Entries))
	}
}

// TestInputCleanup verifies all source files are removed from the input dir.
func TestInputCleanup(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	for _, name := range []string{"a.txt", "b.txt", "c.txt"} {
		if err := os.WriteFile(filepath.Join(inputDir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	runPipeline(t, inputDir, outputDir)

	entries, _ := os.ReadDir(inputDir)
	if len(entries) != 0 {
		t.Errorf("input dir should be empty after processing, got %d files", len(entries))
	}
}

// TestKeyboardNavJS verifies that the generated HTML includes navigation attributes.
func TestKeyboardNavJS(t *testing.T) {
	inputDir, outputDir := makeDirs(t)

	if err := os.WriteFile(filepath.Join(inputDir, "nav.txt"), []byte("nav test"), 0o644); err != nil {
		t.Fatal(err)
	}

	runPipeline(t, inputDir, outputDir)

	index := readFile(t, filepath.Join(outputDir, "index.html"))
	assertContains(t, index, `data-index="0"`, "index.html data-index attribute")
	assertContains(t, index, `.post-active`, "index.html .post-active CSS")
	assertContains(t, index, `playNavSound`, "index.html playNavSound function")
}

// TestThemeSelection verifies that every registered theme renders a valid
// index.html containing core structural elements (post text, nav script).
func TestThemeSelection(t *testing.T) {
	themes := []string{
		"aurora", "brutalist", "cosmos", "matrix", "neon",
		"ocean", "plasma", "retro", "synthwave", "terminal", "volcano",
	}

	for _, theme := range themes {
		theme := theme // capture for parallel sub-test

		t.Run(theme, func(t *testing.T) {
			inputDir, outputDir := makeDirs(t)

			if err := os.WriteFile(filepath.Join(inputDir, "hello.txt"), []byte("theme test post"), 0o644); err != nil {
				t.Fatal(err)
			}

			cfg := &config.Config{
				InputDir:  inputDir,
				OutputDir: outputDir,
				BaseURL:   "https://snonux.foo",
				Theme:     theme,
			}

			if _, err := processor.Run(cfg); err != nil {
				t.Fatalf("processor.Run: %v", err)
			}
			if err := generator.Run(cfg); err != nil {
				t.Fatalf("generator.Run for theme %q: %v", theme, err)
			}

			index := readFile(t, filepath.Join(outputDir, "index.html"))
			assertContains(t, index, "theme test post", "post text")
			assertContains(t, index, "playNavSound", "nav JS")
			assertContains(t, index, `data-index="0"`, "data-index attribute")
		})
	}
}

func sampleRGBAImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 245, B: 255, A: 255})
		}
	}
	return img
}

// writeSamplePNG writes a small 10×10 solid-colour PNG to path.
func writeSamplePNG(t *testing.T, path string) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, sampleRGBAImage()); err != nil {
		t.Fatal(err)
	}
}

// writeSampleJPEG writes a small valid JPEG to path.
func writeSampleJPEG(t *testing.T, path string) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, sampleRGBAImage(), &jpeg.Options{Quality: 90}); err != nil {
		t.Fatal(err)
	}
}

// writeSampleGIF writes a small single-frame GIF to path.
func writeSampleGIF(t *testing.T, path string) {
	t.Helper()

	bounds := image.Rect(0, 0, 10, 10)
	paletted := image.NewPaletted(bounds, palette.Plan9)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			paletted.SetColorIndex(x, y, 1)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if err := gif.Encode(f, paletted, &gif.Options{}); err != nil {
		t.Fatal(err)
	}
}
