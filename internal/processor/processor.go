// Package processor scans the input directory for new source files and converts
// each one into a self-contained post directory under outdir/posts/.
// Supported formats: .txt, .md, .png, .jpg, .jpeg, .gif, .mp3.
// Each processed source file is deleted from the input directory afterward.
package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/post"
)

// Run scans cfg.InputDir and processes every eligible file into a post directory
// under cfg.OutputDir/posts/. Returns the number of posts created.
//
// Images referenced by a .md file in the same input directory are consumed by
// that markdown post and are not processed as independent image posts.
func Run(cfg *config.Config) (int, error) {
	entries, err := os.ReadDir(cfg.InputDir)
	if err != nil {
		return 0, fmt.Errorf("read input dir %s: %w", cfg.InputDir, err)
	}

	postsDir := filepath.Join(cfg.OutputDir, "posts")
	if err := os.MkdirAll(postsDir, 0o755); err != nil {
		return 0, fmt.Errorf("create posts dir: %w", err)
	}

	// Pre-scan markdown files to discover which image filenames they claim.
	// Claimed images are excluded from independent processing.
	claimed := claimedByMarkdown(entries, cfg.InputDir)

	count := 0

	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if claimed[entry.Name()] {
			continue // consumed by a .md post — skip independent processing
		}

		srcPath := filepath.Join(cfg.InputDir, entry.Name())
		if err := processFile(srcPath, postsDir); err != nil {
			return count, fmt.Errorf("process %s: %w", entry.Name(), err)
		}

		count++
	}

	return count, nil
}

// claimedByMarkdown scans all .md entries in inputDir and returns a set of
// image filenames that are referenced within those markdown files.
// Those images should be embedded in the markdown post, not processed alone.
func claimedByMarkdown(entries []os.DirEntry, inputDir string) map[string]bool {
	claimed := make(map[string]bool)

	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".md" {
			continue
		}

		mdPath := filepath.Join(inputDir, entry.Name())
		data, err := os.ReadFile(mdPath)
		if err != nil {
			continue
		}

		for _, imgName := range findLocalImages(string(data), inputDir) {
			claimed[imgName] = true
		}
	}

	return claimed
}

// processFile processes a single input file into a new post directory.
// The source file is removed from the input dir on success.
func processFile(srcPath, postsDir string) error {
	now := time.Now().UTC()
	id := uniqueID(postsDir, now)

	postDir := filepath.Join(postsDir, id)
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		return fmt.Errorf("create post dir %s: %w", id, err)
	}

	p, inboxExtras, err := buildPost(srcPath, postDir, id)
	if err != nil {
		// Clean up the half-created directory to avoid partial state.
		_ = os.RemoveAll(postDir)
		return err
	}

	if err := p.Save(postDir); err != nil {
		_ = os.RemoveAll(postDir)
		return err
	}

	// Remove markdown-referenced inbox images only after the post is persisted
	// (same ordering as the main source file below).
	for _, path := range inboxExtras {
		_ = os.Remove(path)
	}

	// Delete the source file only after the post has been successfully persisted.
	return os.Remove(srcPath)
}

// buildPost dispatches to the appropriate sub-processor based on file extension
// and returns a populated Post ready to be saved. inboxExtras lists absolute
// paths under the input directory to remove after Save succeeds (markdown-local
// images only); other post types return a nil slice.
func buildPost(srcPath, postDir, id string) (*post.Post, []string, error) {
	ext := strings.ToLower(filepath.Ext(srcPath))

	switch ext {
	case ".txt":
		p, err := buildTextPost(srcPath, id)
		return p, nil, err

	case ".md":
		return buildMarkdownPost(srcPath, postDir, id)

	case ".png", ".jpg", ".jpeg", ".gif":
		p, err := buildImagePost(srcPath, postDir, id)
		return p, nil, err

	case ".mp3":
		p, err := buildAudioPost(srcPath, postDir, id)
		return p, nil, err

	default:
		return nil, nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

func buildTextPost(srcPath, id string) (*post.Post, error) {
	html, err := processTxt(srcPath)
	if err != nil {
		return nil, err
	}

	return &post.Post{
		ID:        id,
		Timestamp: time.Now().UTC(),
		PostType:  post.TypeText,
		Content:   html,
	}, nil
}

func buildMarkdownPost(srcPath, postDir, id string) (*post.Post, []string, error) {
	html, localImages, err := processMd(srcPath)
	if err != nil {
		return nil, nil, err
	}

	sourceDir := filepath.Dir(srcPath)

	assets, err := copyLocalImages(localImages, sourceDir, postDir)
	if err != nil {
		return nil, nil, err
	}

	inboxExtras := make([]string, 0, len(localImages))
	for _, name := range localImages {
		inboxExtras = append(inboxExtras, filepath.Join(sourceDir, name))
	}

	return &post.Post{
		ID:        id,
		Timestamp: time.Now().UTC(),
		PostType:  post.TypeMarkdown,
		Content:   html,
		Assets:    assets,
	}, inboxExtras, nil
}

func buildImagePost(srcPath, postDir, id string) (*post.Post, error) {
	filename, html, err := processImage(srcPath, postDir, id)
	if err != nil {
		return nil, err
	}

	return &post.Post{
		ID:        id,
		Timestamp: time.Now().UTC(),
		PostType:  post.TypeImage,
		Content:   html,
		Assets:    []string{filename},
	}, nil
}

func buildAudioPost(srcPath, postDir, id string) (*post.Post, error) {
	filename, html, err := processAudio(srcPath, postDir, id)
	if err != nil {
		return nil, err
	}

	return &post.Post{
		ID:        id,
		Timestamp: time.Now().UTC(),
		PostType:  post.TypeAudio,
		Content:   html,
		Assets:    []string{filename},
	}, nil
}

// copyLocalImages copies referenced image files from sourceDir into postDir.
// Returns the list of filenames that were successfully copied.
func copyLocalImages(filenames []string, sourceDir, postDir string) ([]string, error) {
	var copied []string

	for _, name := range filenames {
		src := filepath.Join(sourceDir, name)
		dst := filepath.Join(postDir, name)

		if err := copyFile(src, dst); err != nil {
			return nil, fmt.Errorf("copy image asset %s: %w", name, err)
		}

		copied = append(copied, name)
	}

	return copied, nil
}

// uniqueID generates a post ID for the given time that does not already exist
// as a directory under postsDir. Appends a numeric suffix if needed.
func uniqueID(postsDir string, t time.Time) string {
	for i := 0; ; i++ {
		id := post.NewID(t, i)
		if _, err := os.Stat(filepath.Join(postsDir, id)); os.IsNotExist(err) {
			return id
		}
	}
}
