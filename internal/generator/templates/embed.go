// Package templates exposes the embedded HTML shell, shared CSS/JS bundles,
// and per-theme assets used by the generator. Themes live as directories
// under templates/themes/<name>/ containing theme.css, theme.js, and meta.json.
// All assets are compiled into the binary via //go:embed so snonux still ships
// as a single self-contained executable.
package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
)

//go:embed shell.tmpl shared/*.tmpl shared/shared.css shared/shared.js themes/*/theme.css themes/*/theme.js themes/*/meta.json
var FS embed.FS

// Shell returns the body of shell.tmpl — the single page template used for
// every generated HTML page. Theme-specific markup is injected at gen time
// from the default theme's meta.json (and at runtime by shared.js for any
// other selected theme).
func Shell() (string, error) {
	b, err := FS.ReadFile("shell.tmpl")
	if err != nil {
		return "", fmt.Errorf("read shell.tmpl: %w", err)
	}
	return string(b), nil
}

// Shared reads a named shared sub-template (currently only "nav").
func Shared(name string) (string, error) {
	b, err := FS.ReadFile(path.Join("shared", name+".tmpl"))
	if err != nil {
		return "", fmt.Errorf("read shared template %q: %w", name, err)
	}
	return string(b), nil
}

// SharedCSS returns the bundled CSS that every page links via shared.css.
// Used by the generator to write dist/shared.css.
func SharedCSS() ([]byte, error) {
	return FS.ReadFile("shared/shared.css")
}

// SharedJS returns the bundled JS that every page references via shared.js.
// Used by the generator to write dist/shared.js.
func SharedJS() ([]byte, error) {
	return FS.ReadFile("shared/shared.js")
}

// ThemeCSS returns the per-theme stylesheet bytes for the named theme.
func ThemeCSS(name string) ([]byte, error) {
	return FS.ReadFile(path.Join("themes", name, "theme.css"))
}

// ThemeJS returns the per-theme script bytes for the named theme.
func ThemeJS(name string) ([]byte, error) {
	return FS.ReadFile(path.Join("themes", name, "theme.js"))
}

// ThemeMeta returns the per-theme meta.json bytes for the named theme.
func ThemeMeta(name string) ([]byte, error) {
	return FS.ReadFile(path.Join("themes", name, "meta.json"))
}

// ThemeNames returns a sorted list of available theme names derived from the
// directories present under templates/themes/.
func ThemeNames() ([]string, error) {
	entries, err := fs.ReadDir(FS, "themes")
	if err != nil {
		return nil, fmt.Errorf("list theme dir: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		names = append(names, e.Name())
	}

	sort.Strings(names)
	return names, nil
}
