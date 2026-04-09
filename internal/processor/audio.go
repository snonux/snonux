package processor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// processAudio copies an .mp3 file into destDir and returns an HTML <audio> snippet.
// The audio element has controls enabled so visitors can play it inline.
func processAudio(srcPath, destDir, postID string) (filename, htmlContent string, err error) {
	outName := filepath.Base(srcPath)
	outPath := filepath.Join(destDir, outName)

	if err := copyFile(srcPath, outPath); err != nil {
		return "", "", err
	}

	// The src attribute is relative to the site root.
	src := fmt.Sprintf("posts/%s/%s", postID, outName)
	html := fmt.Sprintf(
		`<audio controls class="post-audio"><source src="%s" type="audio/mpeg">Your browser does not support audio.</audio>`,
		src,
	)

	return outName, html, nil
}

// copyFile copies the file at src to dst, creating dst if it does not exist.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source %s: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dest %s: %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy %s → %s: %w", src, dst, err)
	}

	return nil
}
