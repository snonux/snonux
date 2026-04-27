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

	// Sync, when true, rsyncs OutputDir to mirror hosts after a successful run.
	Sync bool

	// SyncTargets are the remote hostnames to rsync to when Sync is true.
	// Defaults to ["pi0.lan.buetow.org", "pi1.lan.buetow.org"].
	// Override with SNONUX_SYNC_TARGETS env var (comma-separated).
	SyncTargets []string

	// SyncRemoteDir is the destination directory on each target host.
	// Defaults to "/var/www/html/snonux/". Override with SNONUX_SYNC_REMOTE_DIR env var.
	SyncRemoteDir string
}
