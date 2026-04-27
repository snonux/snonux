package main

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/generator"
)

func TestExpandHome(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	got, err := expandHome(filepath.Join("~", "snonux-test-sub"))
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, "snonux-test-sub")
	if got != want {
		t.Fatalf("got %q; want %q", got, want)
	}

	got, err = expandHome("/no/tilde")
	if err != nil || got != "/no/tilde" {
		t.Fatalf("got %q err %v", got, err)
	}
}

func TestEnsureDir_createsAndRejectsFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	sub := filepath.Join(dir, "newdir")
	if err := ensureDir(sub); err != nil {
		t.Fatal(err)
	}
	if st, err := os.Stat(sub); err != nil || !st.IsDir() {
		t.Fatal("not a dir")
	}
	if err := ensureDir(sub); err != nil {
		t.Fatal(err)
	}

	filePath := filepath.Join(dir, "file")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := ensureDir(filePath); err == nil {
		t.Fatal("expected error when path is file")
	}
}

func TestParseFlags_version(t *testing.T) {
	t.Parallel()

	_, mode, err := parseFlags([]string{"-version"})
	if err != nil {
		t.Fatal(err)
	}
	if mode != modeVersion {
		t.Fatalf("mode %v", mode)
	}
}

func TestParseFlags_listThemes(t *testing.T) {
	t.Parallel()

	_, mode, err := parseFlags([]string{"-list-themes"})
	if err != nil {
		t.Fatal(err)
	}
	if mode != modeListThemes {
		t.Fatalf("mode %v", mode)
	}
}

func TestParseFlags_run(t *testing.T) {
	t.Parallel()

	cfg, mode, err := parseFlags([]string{
		"-input", "/tmp/in",
		"-output", "/tmp/out",
		"-theme", "neon",
		"-base-url", "https://t.test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if mode != modeRun {
		t.Fatalf("mode %v", mode)
	}
	if cfg.Theme != "neon" || cfg.BaseURL != "https://t.test" {
		t.Fatalf("cfg %+v", cfg)
	}
}

func TestParseFlags_sync(t *testing.T) {
	t.Parallel()

	cfg, mode, err := parseFlags([]string{
		"-input", "./in",
		"-output", "./out",
		"-theme", "neon",
		"-sync",
	})
	if err != nil {
		t.Fatal(err)
	}
	if mode != modeRun {
		t.Fatalf("mode %v", mode)
	}
	if !cfg.Sync {
		t.Fatal("expected cfg.Sync")
	}
}

func TestResolvePaths(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		InputDir:  filepath.Join("~", "snonux-test-in"),
		OutputDir: filepath.Join("~", "snonux-test-out"),
	}

	if err := resolvePaths(cfg); err != nil {
		t.Fatal(err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.InputDir != filepath.Join(home, "snonux-test-in") {
		t.Fatalf("input dir: got %q", cfg.InputDir)
	}
	if cfg.OutputDir != filepath.Join(home, "snonux-test-out") {
		t.Fatalf("output dir: got %q", cfg.OutputDir)
	}
}

func TestResolveTheme_fixed(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{Theme: "neon"}
	if err := resolveTheme(cfg, rand.New(rand.NewSource(1))); err != nil {
		t.Fatal(err)
	}
	if cfg.Theme != "neon" {
		t.Fatalf("expected fixed theme, got %q", cfg.Theme)
	}
}

func TestResolveTheme_random(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{Theme: "random"}
	if err := resolveTheme(cfg, rand.New(rand.NewSource(42))); err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, n := range generator.ListThemes() {
		names[n] = true
	}
	if !names[cfg.Theme] {
		t.Fatalf("unexpected theme %q", cfg.Theme)
	}
}

// TestResolveTheme_random_deterministic verifies that when a seeded
// *rand.Rand is passed in, the "random" theme resolve predictably.
func TestResolveTheme_random_deterministic(t *testing.T) {
	t.Parallel()

	cfg1 := &config.Config{Theme: "random"}
	if err := resolveTheme(cfg1, rand.New(rand.NewSource(1))); err != nil {
		t.Fatal(err)
	}

	cfg2 := &config.Config{Theme: "random"}
	if err := resolveTheme(cfg2, rand.New(rand.NewSource(1))); err != nil {
		t.Fatal(err)
	}
	if cfg1.Theme != cfg2.Theme {
		t.Fatalf("expected deterministic theme %q, got %q", cfg1.Theme, cfg2.Theme)
	}
}

func TestResolveTheme_nilRng(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{Theme: "random"}
	if err := resolveTheme(cfg, nil); err == nil {
		t.Fatal("expected error for nil rng")
	}
}

func TestValidateDirs(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()
	cfg := &config.Config{InputDir: in, OutputDir: out}
	if err := validateDirs(cfg); err != nil {
		t.Fatal(err)
	}

	badFile := filepath.Join(t.TempDir(), "notadir")
	if err := os.WriteFile(badFile, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg = &config.Config{InputDir: badFile, OutputDir: out}
	if err := validateDirs(cfg); err == nil {
		t.Fatal("expected error when input is file")
	}
}

func TestParseFlags_unknownFlag(t *testing.T) {
	t.Parallel()

	_, _, err := parseFlags([]string{"-not-a-real-flag"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, errParseFlags) {
		t.Fatalf("got %v", err)
	}
}

func TestRun_pipeline(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()
	if err := os.WriteFile(filepath.Join(in, "a.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		InputDir:  in,
		OutputDir: out,
		BaseURL:   "https://pipe.test",
		Theme:     "neon",
	}
	if err := run(cfg); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "snonux") {
		t.Fatal("expected generated html")
	}
}
