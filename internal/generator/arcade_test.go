package generator

import (
	"bytes"
	"image/jpeg"
	"regexp"
	"testing"

	"github.com/snonux/snonuxmicroblog/internal/generator/templates"
)

// gameThumbName is the picture the game overview menu shows for a theme.
const gameThumbName = "game-thumb.jpg"

// blurbsBlock captures the body of the BLURBS object literal in arcade.js,
// blurbEntry each `name: '…'` line inside it, and gameBlurb the blurb a
// game.js declares for itself.
var (
	blurbsBlock = regexp.MustCompile(`(?s)var BLURBS = \{(.*?)\};`)
	blurbEntry  = regexp.MustCompile(`([a-z0-9]+): '((?:[^'\\]|\\.)*)'`)
	gameBlurb   = regexp.MustCompile(`blurb: '((?:[^'\\]|\\.)*)'`)
)

// menuBlurbs returns theme → blurb as listed in arcade.js; nil when the
// block is missing.
func menuBlurbs(arcade string) map[string]string {
	block := blurbsBlock.FindStringSubmatch(arcade)
	if block == nil {
		return nil
	}
	out := make(map[string]string)
	for _, m := range blurbEntry.FindAllStringSubmatch(block[1], -1) {
		out[m[1]] = m[2]
	}
	return out
}

func TestMenuBlurbs_parsesBlurbsBlock(t *testing.T) {
	t.Parallel()

	// An escaped apostrophe does not end a blurb.
	got := menuBlurbs("var BLURBS = {\n  neon: 'Box them in.',\n  dos: 'Don\\'t dig.'\n};")
	if len(got) != 2 || got["neon"] != "Box them in." || got["dos"] != `Don\'t dig.` {
		t.Errorf("menuBlurbs = %v", got)
	}
	if menuBlurbs("var OTHER = {};") != nil {
		t.Error("menuBlurbs should be nil without a BLURBS block")
	}
}

// TestArcadeMenu_describesEveryGame checks the menu against the games: the
// blurb it shows has to be the one the game itself declares (the menu cannot
// read game.js, which is fetched on first play only), and it must not
// describe a game that does not exist.
func TestArcadeMenu_describesEveryGame(t *testing.T) {
	t.Parallel()

	arcade, err := templates.SharedArcadeJS()
	if err != nil {
		t.Fatalf("read arcade.js: %v", err)
	}
	blurbs := menuBlurbs(string(arcade))
	themes := ListThemes()
	if len(blurbs) != len(themes) {
		t.Errorf("arcade.js describes %d games, want %d", len(blurbs), len(themes))
	}
	for _, theme := range themes {
		script := themeGameScript(t, theme)
		if script == nil {
			continue // reported by TestThemeGames_registerUnderOwnName
		}
		m := gameBlurb.FindSubmatch(script)
		if m == nil {
			t.Errorf("theme %q: game.js declares no blurb", theme)
			continue
		}
		if got, want := blurbs[theme], string(m[1]); got != want {
			t.Errorf("theme %q: arcade.js says %q, game.js says %q", theme, got, want)
		}
	}
}

// TestArcadeMenu_everyGameHasAPicture guards the menu's pictures: a missing
// or broken one shows as an empty box on the card. They are written by
// integrationtests/games/thumbs.mjs at 320x180.
func TestArcadeMenu_everyGameHasAPicture(t *testing.T) {
	t.Parallel()

	for _, theme := range ListThemes() {
		data, err := templates.FS.ReadFile("themes/" + theme + "/" + gameThumbName)
		if err != nil {
			t.Errorf("theme %q: no %s (run integrationtests/games/thumbs.mjs)", theme, gameThumbName)
			continue
		}
		cfg, err := jpeg.DecodeConfig(bytes.NewReader(data))
		if err != nil {
			t.Errorf("theme %q: %s is not a JPEG: %v", theme, gameThumbName, err)
			continue
		}
		if cfg.Width != 320 || cfg.Height != 180 {
			t.Errorf("theme %q: %s is %dx%d, want 320x180", theme, gameThumbName, cfg.Width, cfg.Height)
		}
		if len(data) > 40<<10 {
			t.Errorf("theme %q: %s is %d bytes, want at most 40 KiB", theme, gameThumbName, len(data))
		}
	}
}
