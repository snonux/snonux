// Package templates exposes the embedded HTML theme and shared sub-templates
// used by the generator. Themes and nav/HTML fragments live as separate files
// under templates/themes/*.tmpl and templates/shared/*.tmpl so they can be
// edited without recompiling logic; they are compiled into the binary via
// //go:embed so snonux still ships as a single self-contained executable.
package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

//go:embed themes/*.tmpl shared/*.tmpl
var FS embed.FS

// Theme reads the raw HTML template body for a named theme.
// The returned string is the outer page template; shared sub-templates
// (navhints, navmodal, navscript, etc.) are obtained via Shared().
func Theme(name string) (string, error) {
	b, err := FS.ReadFile(path.Join("themes", name+".tmpl"))
	if err != nil {
		return "", fmt.Errorf("read theme %q: %w", name, err)
	}
	return string(b), nil
}

// Shared reads a named shared sub-template file from shared/*.tmpl.
func Shared(name string) (string, error) {
	b, err := FS.ReadFile(path.Join("shared", name+".tmpl"))
	if err != nil {
		return "", fmt.Errorf("read shared template %q: %w", name, err)
	}
	return string(b), nil
}

// ThemeNames returns a sorted list of available theme names derived from the
// files present under templates/themes/.
func ThemeNames() ([]string, error) {
	entries, err := fs.ReadDir(FS, "themes")
	if err != nil {
		return nil, fmt.Errorf("list theme dir: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		n := e.Name()
		if !strings.HasSuffix(n, ".tmpl") {
			continue
		}
		names = append(names, strings.TrimSuffix(n, ".tmpl"))
	}

	sort.Strings(names)
	return names, nil
}
