package generator

import (
	"log"
	"sync"

	"codeberg.org/snonux/snonux/internal/generator/templates"
)

// fallbackThemeName is used when an unknown name is requested for default
// metadata, matching the previous behaviour of the per-theme registry.
const fallbackThemeName = "neon"

// themeSet caches the list of theme names available in the embedded template FS
// so ListThemes does not re-read the directory on every call.
var (
	themeSetCache map[string]struct{}
	themeSetOnce  sync.Once
)

func getThemeSet() map[string]struct{} {
	themeSetOnce.Do(func() {
		names, err := templates.ThemeNames()
		if err != nil {
			// At build time the embed //go:embed directive guarantees the FS is
			// populated, so this should never happen; log and continue with an
			// empty set so callers can fall back cleanly.
			log.Printf("warning: could not enumerate themes from embedded FS: %v", err)
			themeSetCache = map[string]struct{}{}
			return
		}

		out := make(map[string]struct{}, len(names))
		for _, n := range names {
			out[n] = struct{}{}
		}
		themeSetCache = out
	})
	return themeSetCache
}

// validThemeName returns name if it is a known theme, otherwise the fallback.
// Callers use this to coerce CLI input ("--theme random" already resolves
// upstream) so downstream lookups never miss.
func validThemeName(name string) string {
	if _, ok := getThemeSet()[name]; ok {
		return name
	}
	return fallbackThemeName
}

// ListThemes returns a sorted list of all available theme names.
func ListThemes() []string {
	names, err := templates.ThemeNames()
	if err != nil {
		log.Printf("warning: could not list themes from embedded FS: %v", err)
		return nil
	}
	return names
}
