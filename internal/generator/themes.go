package generator

import (
	"log"

	"codeberg.org/snonux/snonux/internal/generator/templates"
)

// fallbackThemeName is returned by getTheme when an unknown name is requested,
// matching the previous behaviour of the hand-maintained themeRegistry map.
const fallbackThemeName = "neon"

// themeSet caches the list of theme names available in the embedded template FS
// so ListThemes and getTheme do not re-read the directory on every call.
var themeSet = loadThemeSet()

func loadThemeSet() map[string]struct{} {
	names, err := templates.ThemeNames()
	if err != nil {
		// At build time the embed //go:embed directive guarantees the FS is
		// populated, so this should never happen; log and continue with an
		// empty set so getTheme() falls back cleanly.
		log.Printf("warning: could not enumerate themes from embedded FS: %v", err)
		return map[string]struct{}{}
	}

	out := make(map[string]struct{}, len(names))
	for _, n := range names {
		out[n] = struct{}{}
	}
	return out
}

// getTheme returns the HTML template body for the given theme name, loading it
// from the embedded template FS. It falls back to the neon theme if the name
// is unknown (preserving previous behaviour of the hand-maintained map).
func getTheme(name string) string {
	if _, ok := themeSet[name]; !ok {
		name = fallbackThemeName
	}

	body, err := templates.Theme(name)
	if err != nil {
		// Last-resort fallback: try neon. If that also fails, return an empty
		// string; template.Parse will then produce a diagnostic error.
		if body, err = templates.Theme(fallbackThemeName); err != nil {
			log.Printf("warning: could not load fallback theme %q: %v", fallbackThemeName, err)
			return ""
		}
	}

	return body
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
