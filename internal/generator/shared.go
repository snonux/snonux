package generator

import (
	"log"

	"codeberg.org/snonux/snonux/internal/generator/templates"
)

// navDefs is the content of templates/shared/nav.tmpl loaded once at startup.
// It defines named sub-templates shared across all themes:
//   - "splashGate" — synchronous script: first child of <body>; sets html.sno-splash-skip when
//     splash should not run (?splash=0, not index path, or Referer from same-site index/pageN).
//   - "navhints"  — keyboard shortcut hint bar HTML
//   - "navSharedCSSInner" — shared CSS (injected inside each theme's <style> in <head>)
//   - "navmodal"  — full-screen expanded-post modal HTML (no <style>; CSS lives in head)
//   - "navscript" — keyboard navigation + Web Audio; splash/nav/modal sounds from themeSoundsJSON (per theme)
//
// Each theme ends its <style> with {{template "navSharedCSSInner"}} then calls
// {{template "splashGate"}}, {{template "navhints" .}}, {{template "navmodal" .}},
// and {{template "navscript" .}} at the appropriate points in its HTML.
var navDefs = loadNavDefs()

func loadNavDefs() string {
	s, err := templates.Shared("nav")
	if err != nil {
		log.Printf("warning: could not load shared nav template: %v", err)
		return ""
	}
	return s
}
