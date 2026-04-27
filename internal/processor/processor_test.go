package processor

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/post"
)

func TestRun_processesTxt(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()
	if err := os.WriteFile(filepath.Join(in, "note.txt"), []byte("Hello world"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{InputDir: in, OutputDir: out, BaseURL: "https://x.test"}
	n, err := Run(cfg)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Fatalf("processed count = %d; want 1", n)
	}

	entries, err := os.ReadDir(filepath.Join(out, "posts"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("posts dir: %v entries=%v", err, entries)
	}
	postDir := filepath.Join(out, "posts", entries[0].Name())
	p, err := post.Load(postDir)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if p.PostType != post.TypeText {
		t.Fatalf("type %v", p.PostType)
	}
	if _, err := os.ReadFile(filepath.Join(in, "note.txt")); !os.IsNotExist(err) {
		t.Fatal("source should be removed")
	}
}

func TestRun_unsupportedExt(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()
	if err := os.WriteFile(filepath.Join(in, "x.bin"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Run(&config.Config{InputDir: in, OutputDir: out, BaseURL: "https://x"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_readInputDirFails(t *testing.T) {
	t.Parallel()

	_, err := Run(&config.Config{InputDir: "/nonexistent/inbox/xyz", OutputDir: t.TempDir(), BaseURL: "https://x"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_png(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()
	pngPath := filepath.Join(in, "shot.png")
	f, err := os.Create(pngPath)
	if err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	n, err := Run(&config.Config{InputDir: in, OutputDir: out, BaseURL: "https://x"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Fatalf("n=%d", n)
	}
}

func TestRun_mp3(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()
	if err := os.WriteFile(filepath.Join(in, "clip.mp3"), []byte("fake-mp3-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}

	n, err := Run(&config.Config{InputDir: in, OutputDir: out, BaseURL: "https://x"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Fatalf("n=%d", n)
	}
}

func TestRun_markdown(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()
	if err := os.WriteFile(filepath.Join(in, "x.md"), []byte("# Hi\n\n**bold**"), 0o644); err != nil {
		t.Fatal(err)
	}

	n, err := Run(&config.Config{InputDir: in, OutputDir: out, BaseURL: "https://x"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Fatalf("n=%d", n)
	}

	entries, _ := os.ReadDir(filepath.Join(out, "posts"))
	p, err := post.Load(filepath.Join(out, "posts", entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	if p.PostType != post.TypeMarkdown {
		t.Fatalf("got %v", p.PostType)
	}
}

func TestUniqueID_new(t *testing.T) {
	t.Parallel()

	postsDir := t.TempDir()
	id, err := uniqueID(postsDir, time.Now().UTC())
	if err != nil {
		t.Fatalf("uniqueID: %v", err)
	}
	if id == "" {
		t.Fatal("expected non-empty id")
	}
}

func TestUniqueID_collision(t *testing.T) {
	t.Parallel()

	postsDir := t.TempDir()
	now := time.Now().UTC()

	// Pre-create the first expected directory so uniqueID must pick the next suffix.
	firstID := post.NewID(now, 0)
	if err := os.MkdirAll(filepath.Join(postsDir, firstID), 0o755); err != nil {
		t.Fatal(err)
	}

	id, err := uniqueID(postsDir, now)
	if err != nil {
		t.Fatalf("uniqueID: %v", err)
	}
	if id == firstID {
		t.Fatalf("expected different id, got %q", id)
	}
}

func TestUniqueID_statError(t *testing.T) {
	t.Parallel()

	// Create a postsDir and remove read permission so Stat fails with
	// a permission error rather than IsNotExist.
	postsDir := t.TempDir()
	if err := os.Chmod(postsDir, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(postsDir, 0o755) // restore for cleanup

	_, err := uniqueID(postsDir, time.Now().UTC())
	if err == nil {
		t.Fatal("expected error when stat fails")
	}
}

func TestRun_markdownWithLocalImage(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()
	pngPath := filepath.Join(in, "embed.png")
	f, err := os.Create(pngPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	md := `![x](embed.png)
text`
	if err := os.WriteFile(filepath.Join(in, "post.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}

	n, err := Run(&config.Config{InputDir: in, OutputDir: out, BaseURL: "https://x"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Fatalf("n=%d", n)
	}

	entries, _ := os.ReadDir(filepath.Join(out, "posts"))
	pdir := filepath.Join(out, "posts", entries[0].Name())
	p, err := post.Load(pdir)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Assets) < 1 {
		t.Fatalf("want assets, got %+v", p.Assets)
	}
	if _, err := os.Stat(filepath.Join(pdir, "embed.png")); err != nil {
		t.Fatal(err)
	}
}

func TestRun_twoMarkdownsClaimingSameImageFails(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()

	// Shared image in the inbox.
	pngPath := filepath.Join(in, "pic.png")
	f, err := os.Create(pngPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	// Two markdown files both reference the same image.
	if err := os.WriteFile(filepath.Join(in, "a.md"), []byte("![a](pic.png)\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(in, "b.md"), []byte("![b](pic.png)\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err = Run(&config.Config{InputDir: in, OutputDir: out, BaseURL: "https://x"})
	if err == nil {
		t.Fatal("expected error when two markdowns claim the same image")
	}
	if !strings.Contains(err.Error(), "pic.png") {
		t.Fatalf("error should mention the conflicting image, got: %v", err)
	}

	// Verify that no post directories were created and no source files deleted.
	entries, _ := os.ReadDir(filepath.Join(out, "posts"))
	if len(entries) != 0 {
		t.Fatalf("expected no posts created, got %d", len(entries))
	}
	for _, name := range []string{"pic.png", "a.md", "b.md"} {
		if _, err := os.Stat(filepath.Join(in, name)); err != nil {
			t.Fatalf("source %s should still exist: %v", name, err)
		}
	}
}

func TestRun_duplicateImageClaimsInSameMarkdownAllowed(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()

	pngPath := filepath.Join(in, "pic.png")
	f, err := os.Create(pngPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()

	// Same markdown references the same image twice — should not be treated as a conflict.
	md := "![first](pic.png)\n![second](pic.png)\n"
	if err := os.WriteFile(filepath.Join(in, "post.md"), []byte(md), 0o644); err != nil {
		t.Fatal(err)
	}

	n, err := Run(&config.Config{InputDir: in, OutputDir: out, BaseURL: "https://x"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if n != 1 {
		t.Fatalf("n=%d; want 1", n)
	}
}
