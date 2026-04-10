package main

import (
	"errors"
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

	in := t.TempDir()
	out := t.TempDir()
	cfg, mode, err := parseFlags([]string{
		"-input", in,
		"-output", out,
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

	in := t.TempDir()
	out := t.TempDir()
	cfg, mode, err := parseFlags([]string{
		"-input", in,
		"-output", out,
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

func TestParseFlags_randomTheme(t *testing.T) {
	t.Parallel()

	in := t.TempDir()
	out := t.TempDir()
	cfg, _, err := parseFlags([]string{"-input", in, "-output", out, "-theme", "random"})
	if err != nil {
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
