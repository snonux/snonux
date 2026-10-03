package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/snonux/snonux/internal/config"
	"github.com/snonux/snonux/internal/generator/templates"
	"github.com/snonux/snonux/internal/post"
)

// gameScriptName is the per-theme file that holds a theme's arcade game.
const gameScriptName = "game.js"

// registerCall matches SnoGame.register('<theme>', …) however the engine
// object is named inside the script (G.register, SnoGame.register, …).
var registerCall = regexp.MustCompile(`\.register\(\s*['"]([^'"]*)['"]`)

// checkGameScript verifies that a theme's game.js registers exactly one game,
// under the theme's own name, and that the engine manifest lists that theme.
// shared/games.js launches a game by theme name, so a mismatch in either
// place would leave the launch button dead.
func checkGameScript(theme string, script []byte, engine string) error {
	calls := registerCall.FindAllSubmatch(script, -1)
	if len(calls) != 1 {
		return fmt.Errorf("theme %q: want exactly 1 register call in %s, found %d", theme, gameScriptName, len(calls))
	}
	if got := string(calls[0][1]); got != theme {
		return fmt.Errorf("theme %q: %s registers as %q", theme, gameScriptName, got)
	}
	if !regexp.MustCompile(`\b` + regexp.QuoteMeta(theme) + `: '`).MatchString(engine) {
		return fmt.Errorf("theme %q: missing from the TITLES manifest in games.js", theme)
	}
	return nil
}

// themeGameScript returns the theme's game.js, or nil when it has none.
func themeGameScript(t *testing.T, theme string) []byte {
	t.Helper()

	extras, err := templates.ThemeExtraFiles(theme)
	if err != nil {
		t.Fatalf("list %s extras: %v", theme, err)
	}
	for _, f := range extras {
		if f.Name == gameScriptName {
			return f.Data
		}
	}
	return nil
}

func TestThemeGames_registerUnderOwnName(t *testing.T) {
	t.Parallel()

	engine, err := templates.SharedGamesJS()
	if err != nil {
		t.Fatalf("read games.js: %v", err)
	}

	found := 0
	for _, theme := range ListThemes() {
		script := themeGameScript(t, theme)
		if script == nil {
			continue
		}
		found++
		if err := checkGameScript(theme, script, string(engine)); err != nil {
			t.Error(err)
		}
	}
	if found == 0 {
		t.Fatal("no theme ships a game.js")
	}
}

func TestCheckGameScript_rejectsBrokenScripts(t *testing.T) {
	t.Parallel()

	engine := "var TITLES = { neon: 'Light Cycles' };"
	cases := []struct {
		name   string
		theme  string
		script string
		want   string
	}{
		{"valid", "neon", "G.register('neon', {});", ""},
		{"empty script", "neon", "", "found 0"},
		{"wrong name", "neon", "G.register('ocean', {});", `registers as "ocean"`},
		{"two games", "neon", "G.register('neon', {}); G.register('neon', {});", "found 2"},
		{"not in manifest", "ocean", `SnoGame.register("ocean", {});`, "missing from the TITLES manifest"},
	}

	for _, tc := range cases {
		err := checkGameScript(tc.theme, []byte(tc.script), engine)
		switch {
		case tc.want == "" && err != nil:
			t.Errorf("%s: unexpected error %v", tc.name, err)
		case tc.want != "" && (err == nil || !strings.Contains(err.Error(), tc.want)):
			t.Errorf("%s: got error %v, want one containing %q", tc.name, err, tc.want)
		}
	}
}

func TestRun_writesGameAssets(t *testing.T) {
	t.Parallel()

	out := t.TempDir()
	postDir := filepath.Join(out, "posts", "g1")
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		t.Fatal(err)
	}
	p := &post.Post{ID: "g1", Timestamp: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC), PostType: post.TypeText, Content: "<p>hi</p>"}
	if err := p.Save(postDir); err != nil {
		t.Fatal(err)
	}
	if err := Run(ctx, &config.Config{OutputDir: out, BaseURL: "https://example.test", Theme: "neon"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	wantIn := map[string]string{
		"games.js":                  "window.SnoGame",
		"index.html":                `<script src="games.js" defer></script>`,
		"shared.css":                "#sno-game",
		"shared.js":                 "snonuxGameLaunch",
		"themes/breakout/game.js":   "G.register('breakout'",
		"themes/breakout/theme.css": "",
	}
	for name, needle := range wantIn {
		data, err := os.ReadFile(filepath.Join(out, filepath.FromSlash(name)))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if !strings.Contains(string(data), needle) {
			t.Errorf("%s does not contain %q", name, needle)
		}
	}

	// The engine is a shared asset: it must not be duplicated into themes.
	if _, err := os.Stat(filepath.Join(out, "themes", "neon", "games.js")); !os.IsNotExist(err) {
		t.Errorf("themes/neon/games.js should not exist (err=%v)", err)
	}
}
