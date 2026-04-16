package generator

// HTML theme templates and shared sub-templates live as separate files under
// internal/generator/templates/{themes,shared}/*.tmpl and are loaded into the
// binary via embed.FS; see internal/generator/templates/embed.go.
//
// themes.go exposes the theme registry and getTheme()/ListThemes() helpers.
// shared.go loads templates/shared/nav.tmpl into navDefs.
// favicon.go loads templates/shared/favicon_head.tmpl into faviconHeadHTML.
