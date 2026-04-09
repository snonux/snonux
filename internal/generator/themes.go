package generator

// themeRegistry maps theme names to their HTML template strings.
// Each template must use {{template "navhints" .}}, {{template "navmodal" .}},
// and {{template "navscript" .}} — these are defined in shared.go (navDefs).
var themeRegistry = map[string]string{
	"neon":      neonTemplate,
	"terminal":  terminalTemplate,
	"synthwave": synthwaveTemplate,
	"minimal":   minimalTemplate,
	"brutalist": brutalistTemplate,
	"paper":     paperTemplate,
	"aurora":    auroraTemplate,
	"matrix":    matrixTemplate,
	"ocean":     oceanTemplate,
	"retro":     retroTemplate,
	"glass":     glassTemplate,
}

// getTheme returns the HTML template string for the given theme name.
// Falls back to the neon theme if the name is unknown.
func getTheme(name string) string {
	if t, ok := themeRegistry[name]; ok {
		return t
	}
	return neonTemplate
}

// ListThemes returns a sorted list of all available theme names.
func ListThemes() []string {
	names := make([]string, 0, len(themeRegistry))
	for k := range themeRegistry {
		names = append(names, k)
	}
	// Sort for deterministic output in --help text.
	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] > names[j] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}
	return names
}
