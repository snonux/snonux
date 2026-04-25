// Package generator builds static HTML pages and delegates Atom feed output to
// subpackage atom.
//
// Responsibilities by area (file → role):
//
//   - generator.go — Orchestration: load posts from disk, sort newest-first,
//     paginate, parse shell+nav templates, write index.html / pageN.html plus
//     shared.css/shared.js and per-theme assets, then call atom.Generate.
//   - themes.go — ListThemes / validThemeName helpers backed by the embedded FS.
//   - shared.go — navDefs: shared {{define}} blocks (splashGate, navhints,
//     navmodal) merged at parse time with shell.tmpl so a single html/template
//     parse sees every name.
//   - theme_sounds.go — Per-theme Web Audio parameters; one file is written per
//     theme to dist/themes/<name>/sounds.json and the default theme's preset is
//     also baked into shell.tmpl as window.SNONUX_SOUNDS so the splash chime
//     can fire instantly.
//   - favicon.go — Generates the 32×32 .ico binary written into the output dir.
//   - templates.go — Short pointer: where templates and helpers live.
//
// Theme assets live as separate files under templates/themes/<name>/theme.css,
// theme.js, and meta.json. The single shell.tmpl loads the active theme's CSS
// synchronously in <head> (chosen from localStorage by a boot script) and
// shared.js orchestrates the rest. Switching themes is a localStorage write
// followed by location.reload().
//
// Dependency direction: shell and shared nav templates are composed only for
// the HTML path (generator.go). Package atom depends on config and post only,
// not on themes or html/template, so feed logic stays isolated from page chrome.
package generator
