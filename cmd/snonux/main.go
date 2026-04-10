// Command snonux is the static microblog generator for snonux.foo.
// It processes new source files from the input directory into post directories,
// then regenerates all HTML pages and the Atom feed in the output directory.
//
// Usage:
//
//	snonux --input ./inbox --output ./outdir [--base-url https://snonux.foo]
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/generator"
	"codeberg.org/snonux/snonux/internal/processor"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := run(cfg); err != nil {
		log.Fatalf("error: %v", err)
	}
}

// parseFlags reads CLI flags and returns a validated Config.
// Special theme value "random" picks a theme at random from the registry.
func parseFlags() (*config.Config, error) {
	cfg := &config.Config{}
	listThemes := flag.Bool("list-themes", false, "print all available theme names and exit")

	flag.StringVar(&cfg.InputDir, "input", "./inbox", "directory containing new source files to process")
	flag.StringVar(&cfg.OutputDir, "output", "~/git/snonux.foo/dist", "root directory for generated static site output")
	flag.StringVar(&cfg.BaseURL, "base-url", "https://snonux.foo", "canonical base URL used in Atom feed links")
	flag.StringVar(&cfg.Theme, "theme", "neon", "visual theme name, or \"random\" to pick one at random")
	flag.Parse()

	if *listThemes {
		fmt.Println(strings.Join(generator.ListThemes(), "\n"))
		os.Exit(0)
	}

	// Resolve the special "random" value before any further validation.
	if cfg.Theme == "random" {
		themes := generator.ListThemes()
		cfg.Theme = themes[rand.Intn(len(themes))]
		log.Printf("random theme selected: %s", cfg.Theme)
	}

	var err error

	cfg.InputDir, err = expandHome(cfg.InputDir)
	if err != nil {
		return nil, fmt.Errorf("input dir: %w", err)
	}

	cfg.OutputDir, err = expandHome(cfg.OutputDir)
	if err != nil {
		return nil, fmt.Errorf("output dir: %w", err)
	}

	if err := ensureDir(cfg.InputDir); err != nil {
		return nil, fmt.Errorf("input dir: %w", err)
	}

	if err := ensureDir(cfg.OutputDir); err != nil {
		return nil, fmt.Errorf("output dir: %w", err)
	}

	return cfg, nil
}

// expandHome replaces a leading ~ with the current user's home directory.
func expandHome(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}

	return filepath.Join(home, path[1:]), nil
}

// run executes both pipeline phases: process inputs, then regenerate pages.
func run(cfg *config.Config) error {
	processed, err := processor.Run(cfg)
	if err != nil {
		return fmt.Errorf("processing input files: %w", err)
	}

	log.Printf("processed %d new post(s) from %s", processed, cfg.InputDir)

	if err := generator.Run(cfg); err != nil {
		return fmt.Errorf("generating site: %w", err)
	}

	log.Printf("site regenerated in %s", cfg.OutputDir)

	return nil
}

// ensureDir creates dir if it does not exist, or returns an error if path
// exists but is not a directory.
func ensureDir(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0o755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", dir)
	}

	return nil
}
