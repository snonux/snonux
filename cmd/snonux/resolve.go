package main

import (
	"fmt"
	"math/rand"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/generator"
)

// resolvePaths expands home directories in cfg.InputDir and cfg.OutputDir.
func resolvePaths(cfg *config.Config) error {
	var err error

	cfg.InputDir, err = expandHome(cfg.InputDir)
	if err != nil {
		return fmt.Errorf("input dir: %w", err)
	}

	cfg.OutputDir, err = expandHome(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("output dir: %w", err)
	}

	return nil
}

// validateDirs ensures cfg.InputDir and cfg.OutputDir exist as directories.
func validateDirs(cfg *config.Config) error {
	if err := ensureDir(cfg.InputDir); err != nil {
		return fmt.Errorf("input dir: %w", err)
	}
	if err := ensureDir(cfg.OutputDir); err != nil {
		return fmt.Errorf("output dir: %w", err)
	}
	return nil
}

// resolveTheme resolves the special "random" theme value by picking a registered
// theme using rng. The rng parameter must be non-nil.
func resolveTheme(cfg *config.Config, rng *rand.Rand) error {
	if cfg.Theme != "random" {
		return nil
	}
	if rng == nil {
		return fmt.Errorf("theme %q requires a seeded rng", cfg.Theme)
	}
	themes := generator.ListThemes()
	cfg.Theme = themes[rng.Intn(len(themes))]
	return nil
}
