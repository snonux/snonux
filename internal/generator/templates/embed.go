// Package templates exposes the embedded HTML shell, shared CSS/JS bundles,
// and per-theme assets used by the generator. Themes live as directories
// under templates/themes/<name>/ containing theme.css, theme.js, meta.json,
// and sounds.json.
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

// Per-theme directory convention (under templates/themes/<name>/):
//
//   theme.css           — required, always present
//   theme.js            — required, always present
//   meta.json           — required, always present
//   sounds.json         — required, always present
//   <Family-Weight>.woff2 — optional, one or more self-hosted web fonts
//   <Family-Weight>.woff  — optional, fallback / non-woff2 web fonts
//   FONT_LICENSE.txt    — required iff any font file is present;
//                         contains attribution + license + source URL
//
// All fonts are served from the same origin as the site (no third-party
// CDNs). The four required files above have dedicated accessors; every
// other file in a theme directory is shipped verbatim by ThemeExtraFiles
// and writeThemeAsset.
//
// NOTE on the embed glob below: //go:embed requires every pattern to
// match at least one file, so each font extension is only added once
// the first theme actually ships such a file. When you introduce the
// first .woff2 file (or any new extension), append the matching glob
// here, e.g.:  themes/*/*.woff2
//
//go:embed shell.tmpl shared/*.tmpl shared/shared.css shared/shared.js themes/*/theme.css themes/*/theme.js themes/*/meta.json themes/*/sounds.json themes/*/*.woff themes/*/FONT_LICENSE.txt
var FS embed.FS

// themeStandardFiles lists the per-theme files that have dedicated
// accessor functions and that writeThemeAsset writes by name.
// Any other file inside a theme directory (fonts, license notices, …)
// is treated as an "extra" asset and copied verbatim by ThemeExtraFiles.
var themeStandardFiles = map[string]struct{}{
	"theme.css":   {},
	"theme.js":    {},
	"meta.json":   {},
	"sounds.json": {},
}

// ThemeExtraFile represents a per-theme asset (e.g. a .woff font or a
// FONT_LICENSE.txt) that is not one of the four standard files but should
// still be copied verbatim into dist/themes/<name>/.
type ThemeExtraFile struct {
	// Name is the basename of the file inside the theme directory.
	Name string
	// Data is the raw file content.
	Data []byte
}

// ThemeExtraFiles returns every file under templates/themes/<name>/ that
// is NOT one of the four standard files (theme.css, theme.js, meta.json,
// sounds.json). The returned slice is sorted by name for deterministic
// output.
func ThemeExtraFiles(name string) ([]ThemeExtraFile, error) {
	dir := path.Join("themes", name)
	entries, err := fs.ReadDir(FS, dir)
	if err != nil {
		return nil, fmt.Errorf("list theme dir %q: %w", name, err)
	}

	out := make([]ThemeExtraFile, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if _, ok := themeStandardFiles[e.Name()]; ok {
			continue
		}
		b, err := FS.ReadFile(path.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read theme extra %q/%q: %w", name, e.Name(), err)
		}
		out = append(out, ThemeExtraFile{Name: e.Name(), Data: b})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

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

// ThemeSounds returns the per-theme sounds.json bytes for the named theme.
func ThemeSounds(name string) ([]byte, error) {
	return FS.ReadFile(path.Join("themes", name, "sounds.json"))
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
