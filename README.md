# snonux — Static microblog generator

snonux is the static microblog engine behind [snonux.foo](https://snonux.foo). It processes source files from an input directory into a self-contained static site with paginated HTML pages, an Atom feed, and themed WebGL backgrounds.

## Quick start

```sh
go build -o snonux ./cmd/snonux
./snonux --input ./inbox --output ./dist
```

## Creating posts

Drop files into the input directory (`./inbox` by default). Each file becomes one post. Supported formats:

### Plain text (.txt)

```
inbox/
  thoughts.txt
```

The text content is rendered as-is into a post.

### Markdown (.md)

```
inbox/
  update.md
```

Standard Markdown is converted to HTML (GitHub Flavored Markdown supported). Raw HTML blocks are passed through.

### Markdown with embedded images

Reference a local image from your Markdown file using standard `![alt](filename)` syntax. Place the image file in the same input directory:

```
inbox/
  update.md
  screenshot.png
```

Where `update.md` contains:

```markdown
Check out this screenshot!

![screenshot](screenshot.png)

Pretty neat, right?
```

The image file is automatically copied into the post's asset directory, and the `<img>` src is rewritten to the correct path. The image file is consumed together with the Markdown file and removed from the input directory after processing.

**Note:** The image filename in the Markdown must match the actual file in the inbox. Remote URLs (`http://`, `https://`) are left as-is and not downloaded.

### Images (.png, .jpg, .gif)

```
inbox/
  photo.jpg
```

A standalone image file becomes its own post. Images wider than 1024px are downscaled and re-encoded as JPEG at 80% quality.

### Audio (.mp3)

```
inbox/
  voice-note.mp3
```

An audio file becomes a post with an embedded HTML5 audio player.

### After processing

All source files are removed from the input directory once they have been successfully processed into the output directory.

## Command-line flags

```
--input DIR      Input directory for new source files (default: ./inbox)
--output DIR     Output directory for generated site (default: ./dist)
--base-url URL   Base URL for Atom feed links (default: https://snonux.foo)
--theme NAME     Visual theme, or "random" (default: random)
--sync           Rsync output to pi0/pi1 after generation
--list-themes    Print available theme names and exit
--version        Print version and exit
```

## Themes

Each run can use a different visual theme. Use `--list-themes` to see all available themes, or `--theme random` (the default) to pick one at random.

## Output structure

```
dist/
  index.html                        # Page 1 (newest posts)
  page2.html                        # Page 2, etc.
  atom.xml                          # Atom feed (last 42 entries)
  favicon.ico
  posts/
    2026-04-16-120000/
      post.json                     # Post metadata and rendered HTML
      screenshot.png                # Asset (if any)
    ...
```

## License

See [LICENSE](LICENSE).

### Third-party assets

snonux self-hosts every third-party asset it ships — there are no
runtime requests to Google Fonts, gstatic, int10h, or any other CDN.
Each font file lives next to its theme stylesheet under
`internal/generator/templates/themes/<name>/`, is committed to git,
and is shipped to `dist/themes/<name>/` together with a
`FONT_LICENSE.txt` containing the source URL, version/date, and full
attribution.

Bundled web fonts:

- **dos** — *WebPlus IBM VGA 8x16* (.woff) by VileR, from the
  [Ultimate Oldschool PC Font Pack v2.2](https://int10h.org/oldschool-pc-fonts/),
  [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
  See [internal/generator/templates/themes/dos/FONT_LICENSE.txt](internal/generator/templates/themes/dos/FONT_LICENSE.txt).
- **biomech** — *Oxanium* Regular + Bold (.woff2, latin+latin-ext) by
  Tyler Finck, from [Google Fonts](https://fonts.google.com/specimen/Oxanium),
  [SIL OFL 1.1](https://openfontlicense.org/open-font-license-official-text/).
  See [internal/generator/templates/themes/biomech/FONT_LICENSE.txt](internal/generator/templates/themes/biomech/FONT_LICENSE.txt).

When adding a new bundled font:

1. Download the actual file (`.woff2` preferred, `.woff` accepted) into
   `internal/generator/templates/themes/<name>/` and commit it.
2. Add or update `internal/generator/templates/themes/<name>/FONT_LICENSE.txt`
   with attribution, license, source URL, and version/date.
3. Add an `@font-face` block in that theme's `theme.css` using a
   relative `url(...)` to the bundled file.
4. Add a bullet here pointing at the new `FONT_LICENSE.txt`.

The embed pipeline ([internal/generator/templates/embed.go](internal/generator/templates/embed.go))
automatically picks up `*.woff`, `*.woff2`, and `FONT_LICENSE.txt`
inside any theme directory, and `writeThemeAsset` copies them out to
`dist/themes/<name>/` — no Go code changes needed for new fonts.
