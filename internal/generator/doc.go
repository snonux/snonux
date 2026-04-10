// Package generator builds static HTML pages and delegates Atom feed output to
// subpackage atom.
//
// Responsibilities by area (file → role):
//
//   - generator.go — Orchestration: load posts from disk, sort newest-first,
//     paginate, parse theme+nav templates, write index.html / pageN.html,
//     then call atom.Generate.
//   - themes.go — Theme registry (name → template string) and getTheme /
//     ListThemes for the CLI.
//   - theme_*.go — One file per visual theme: full-page HTML that invokes
//     {{template "navhints" .}}, {{template "navmodal" .}}, {{template "navscript" .}}.
//   - shared.go — navDefs: shared {{define}} blocks merged at parse time with
//     the chosen theme so a single html/template parse sees every name.
//   - templates.go — Short pointer: where templates and registry live.
//
// Dependency direction: themes and shared nav templates are composed only for
// the HTML path (generator.go). Package atom depends on config and post only,
// not on themes or html/template, so feed logic stays isolated from page chrome.
package generator
