package generator

// The unified shell template, shared CSS/JS bundles, and per-theme asset
// triples live under internal/generator/templates/ and are loaded into the
// binary via embed.FS; see internal/generator/templates/embed.go.
//
// themes.go exposes ListThemes()/validThemeName().
// shared.go loads templates/shared/nav.tmpl into getNavDefs() (splashGate, navhints,
// navmodal partials called from shell.tmpl).
// favicon.go generates the favicon.ico binary written into each output dir.
