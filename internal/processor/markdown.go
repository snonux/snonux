package processor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// imageRefPattern matches Markdown image syntax: ![alt](filename)
// We use it to discover local asset references that must be copied.
var imageRefPattern = regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)

// processMd converts a Markdown file to an HTML snippet for a trusted inbox source.
// The markdown (including any raw HTML blocks) is treated as author-controlled
// content, not user-generated input from strangers; see the package comment.
//
// Returns the HTML and a list of local image filenames referenced in the document.
// Referenced images that exist alongside the source file are returned so the
// caller can copy them into the post asset directory.
func processMd(path string) (htmlContent string, localImages []string, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("read markdown %s: %w", path, err)
	}

	// Collect local image references so the caller can copy them as assets.
	localImages = findLocalImages(string(data), filepath.Dir(path))

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			// Trusted inbox: preserve raw HTML in markdown (see package comment).
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(data, &buf); err != nil {
		return "", nil, fmt.Errorf("convert markdown %s: %w", path, err)
	}

	return buf.String(), localImages, nil
}

// findLocalImages returns image filenames referenced in markdown that actually
// exist in sourceDir. Remote URLs (http/https) are ignored.
func findLocalImages(mdContent, sourceDir string) []string {
	matches := imageRefPattern.FindAllStringSubmatch(mdContent, -1)
	var locals []string

	for _, m := range matches {
		ref := m[1]
		// Skip remote URLs.
		if strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://") {
			continue
		}

		candidate := filepath.Join(sourceDir, ref)
		if _, err := os.Stat(candidate); err == nil {
			locals = append(locals, filepath.Base(ref))
		}
	}

	return locals
}
