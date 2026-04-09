// Package config holds the shared configuration for the snonux generator.
// All values are derived from CLI flags with sensible defaults.
package config

const (
	// PostsPerPage is the maximum number of blog posts rendered on a single HTML page.
	PostsPerPage = 42
)

// Config carries the runtime configuration for the generator pipeline.
type Config struct {
	// InputDir is where new source files (txt, md, images, audio) are read from.
	InputDir string

	// OutputDir is the root of the static site: index.html, pageN.html, atom.xml,
	// and the posts/ subdirectory all live here.
	OutputDir string

	// BaseURL is the canonical site URL, used in the Atom feed links.
	// Example: "https://snonux.foo"
	BaseURL string

	// Theme selects the visual style for generated HTML pages.
	// Defaults to "neon". Run with --help to see all available themes.
	Theme string
}
