package processor

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

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
