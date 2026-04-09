package processor

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"
)

// urlPattern matches http/https URLs in plain text.
// Trailing sentence punctuation is stripped separately by stripURLTrailing.
var urlPattern = regexp.MustCompile(`https?://\S+`)

// processTxt reads a plain-text file and wraps each non-empty paragraph in <p> tags.
// URLs are automatically converted to clickable <a> links.
// Non-URL text is HTML-escaped to prevent XSS.
func processTxt(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read txt %s: %w", path, err)
	}

	raw := strings.TrimSpace(string(data))
	if raw == "" {
		return "<p></p>", nil
	}

	// Split on blank lines to get logical paragraphs.
	paragraphs := strings.Split(raw, "\n\n")
	var sb strings.Builder

	for _, para := range paragraphs {
		trimmed := strings.TrimSpace(para)
		if trimmed == "" {
			continue
		}
		fmt.Fprintf(&sb, "<p>%s</p>\n", formatParagraph(trimmed))
	}

	return sb.String(), nil
}

// formatParagraph formats a single paragraph: auto-links URLs, escapes non-URL
// text, and converts single newlines to <br> line breaks.
func formatParagraph(para string) string {
	lines := strings.Split(para, "\n")
	formatted := make([]string, 0, len(lines))

	for _, line := range lines {
		if t := strings.TrimSpace(line); t != "" {
			formatted = append(formatted, autolinkLine(t))
		}
	}

	return strings.Join(formatted, "<br>\n")
}

// autolinkLine escapes non-URL text and wraps detected URLs in <a> tags.
// Opens in a new tab with rel="noopener noreferrer" for security.
func autolinkLine(line string) string {
	locs := urlPattern.FindAllStringIndex(line, -1)
	if len(locs) == 0 {
		return html.EscapeString(line)
	}

	var sb strings.Builder
	prev := 0

	for _, loc := range locs {
		sb.WriteString(html.EscapeString(line[prev:loc[0]]))

		rawURL := line[loc[0]:loc[1]]
		cleanURL := stripURLTrailing(rawURL)
		trailing := rawURL[len(cleanURL):]

		fmt.Fprintf(&sb, `<a href="%s" target="_blank" rel="noopener noreferrer">%s</a>`,
			html.EscapeString(cleanURL), html.EscapeString(cleanURL))

		if trailing != "" {
			sb.WriteString(html.EscapeString(trailing))
		}

		prev = loc[1]
	}

	sb.WriteString(html.EscapeString(line[prev:]))

	return sb.String()
}

// stripURLTrailing removes common sentence-ending punctuation from the end of a
// URL match. These characters are valid in URLs but almost never appear there
// at the end in prose (e.g. "Visit https://foo.com." — the "." ends the sentence).
func stripURLTrailing(u string) string {
	const cutset = ".,;:!?\"')>]}"

	for len(u) > 0 && strings.ContainsRune(cutset, rune(u[len(u)-1])) {
		u = u[:len(u)-1]
	}

	return u
}
