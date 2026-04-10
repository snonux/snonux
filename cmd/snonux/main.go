// Command snonux is the static microblog generator for snonux.foo.
// It processes new source files from the input directory into post directories,
// then regenerates all HTML pages and the Atom feed in the output directory.
//
// Usage:
//
//	snonux [--input ./inbox] [--output ./dist] [--base-url https://snonux.foo]
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/generator"
	"codeberg.org/snonux/snonux/internal/processor"
	"codeberg.org/snonux/snonux/internal/version"
)

// cliMode tells main whether to run the pipeline or print and exit.
type cliMode int

const (
	modeRun cliMode = iota
	modeVersion
	modeListThemes
)

func main() {
	cfg, mode, err := parseFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	switch mode {
	case modeVersion:
		fmt.Println(version.Version)
		return
	case modeListThemes:
		fmt.Println(strings.Join(generator.ListThemes(), "\n"))
		return
	}

	if err := run(cfg); err != nil {
		log.Fatalf("error: %v", err)
	}

	if cfg.Sync {
		if err := syncOutput(cfg.OutputDir); err != nil {
			log.Fatalf("error: %v", err)
		}
	}
}

// errParseFlags is returned when flag parsing fails (e.g. unknown flag).
var errParseFlags = errors.New("flag parse error")

// parseFlags reads CLI flags and returns a validated Config.
// Special theme value "random" picks a theme at random from the registry.
func parseFlags(args []string) (*config.Config, cliMode, error) {
	cfg := &config.Config{}
	fs := flag.NewFlagSet("snonux", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var showVersion bool
	fs.BoolVar(&showVersion, "version", false, "print version and exit (-version, --version)")
	fs.BoolVar(&showVersion, "v", false, "print version and exit (shorthand for -version)")
	listThemes := fs.Bool("list-themes", false, "print all available theme names and exit")

	fs.StringVar(&cfg.InputDir, "input", "./inbox", "directory containing new source files to process")
	fs.StringVar(&cfg.OutputDir, "output", "./dist", "root directory for generated static site output")
	fs.StringVar(&cfg.BaseURL, "base-url", "https://snonux.foo", "canonical base URL used in Atom feed links")
	fs.StringVar(&cfg.Theme, "theme", "random", "visual theme name, or \"random\" to pick one at random")
	fs.BoolVar(&cfg.Sync, "sync", false, "after a successful run, rsync -output to pi0/pi1 when both are pingable (SSH user: SNONUX_SYNC_USER or login name)")

	if err := fs.Parse(args); err != nil {
		return nil, modeRun, fmt.Errorf("%w: %w", errParseFlags, err)
	}

	if showVersion {
		return nil, modeVersion, nil
	}

	if *listThemes {
		return nil, modeListThemes, nil
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
		return nil, modeRun, fmt.Errorf("input dir: %w", err)
	}

	cfg.OutputDir, err = expandHome(cfg.OutputDir)
	if err != nil {
		return nil, modeRun, fmt.Errorf("output dir: %w", err)
	}

	if err := ensureDir(cfg.InputDir); err != nil {
		return nil, modeRun, fmt.Errorf("input dir: %w", err)
	}

	if err := ensureDir(cfg.OutputDir); err != nil {
		return nil, modeRun, fmt.Errorf("output dir: %w", err)
	}

	return cfg, modeRun, nil
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
