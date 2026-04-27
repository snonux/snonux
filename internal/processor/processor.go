// Package processor scans the input directory for new source files and converts
// each one into a self-contained post directory under outdir/posts/.
// Supported formats: .txt, .md, .png, .jpg, .jpeg, .gif, .mp3.
// Each processed source file is deleted from the input directory afterward.
//
// Processing uses a two-phase commit pattern:
//   1. Scan and validate every inbox item without mutating anything.
//   2. Only after all items pass validation, execute mutations
//      (create directories, write assets, persist posts, remove sources).
// If validation fails for any item, the entire batch is aborted and the inbox
// is left untouched. If a mutation fails mid-batch, earlier items have already
// been committed; the failing item is rolled back and the error is returned
// together with the count of successfully committed posts.
//
// Markdown trust boundary: .md files are expected only from a trusted personal
// inbox (the operator’s own email or equivalent). Goldmark is configured with
// html.WithUnsafe so raw HTML and GFM features in those files pass through to
// post HTML intentionally. This is not a multi-tenant or public-submission
// pipeline; do not point an untrusted drop folder at the same input directory
// without replacing that rendering path with sanitization or a stricter parser.
package processor

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"codeberg.org/snonux/snonux/internal/config"
	"codeberg.org/snonux/snonux/internal/post"
)

// Run scans cfg.InputDir and processes every eligible file into a post directory
// under cfg.OutputDir/posts/. It uses a two-phase commit pattern:
//
//   Phase 1 — scan and validate all inbox items without mutating anything.
//   Phase 2 — only after all items pass validation, execute mutations
//             (create directories, write assets, persist posts, remove sources).
//
// If Phase 1 fails for any item, no mutations occur and the inbox is left untouched.
// Returns the number of posts successfully created in this invocation.
func Run(cfg *config.Config) (int, error) {
	entries, err := os.ReadDir(cfg.InputDir)
	if err != nil {
		return 0, fmt.Errorf("read input dir %s: %w", cfg.InputDir, err)
	}

	postsDir := filepath.Join(cfg.OutputDir, "posts")
	if err := os.MkdirAll(postsDir, 0o755); err != nil {
		return 0, fmt.Errorf("create posts dir: %w", err)
	}

	claimed, err := claimedByMarkdown(entries, cfg.InputDir)
	if err != nil {
		return 0, err
	}

	// Phase 1 — validate everything, collect work, mutate nothing.
	var plans []postPlan
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if claimed[entry.Name()] {
			continue
		}

		srcPath := filepath.Join(cfg.InputDir, entry.Name())
		plan, err := planPost(srcPath)
		if err != nil {
			return 0, fmt.Errorf("plan %s: %w", entry.Name(), err)
		}
		plans = append(plans, plan)
	}

	// Phase 2 — commit all mutations.
	count := 0
	now := time.Now().UTC()
	for _, plan := range plans {
		if err := commitPlan(plan, postsDir, now); err != nil {
			return count, fmt.Errorf("commit %s: %w", filepath.Base(plan.srcPath), err)
		}
		count++
	}

	return count, nil
}

// postPlan captures everything validated in Phase 1 for a single source file.
// No file-system mutations are recorded here; only validated content.
type postPlan struct {
	srcPath        string
	ext            string
	textHTML       string
	mdHTML         string
	localImages    []string
	validatedImage image.Image
}

// planPost validates a single source file and returns a plan containing
// everything needed to commit it later. It performs no mutations.
func planPost(srcPath string) (postPlan, error) {
	ext := strings.ToLower(filepath.Ext(srcPath))
	plan := postPlan{srcPath: srcPath, ext: ext}

	switch ext {
	case ".txt":
		html, err := processTxt(srcPath)
		if err != nil {
			return postPlan{}, err
		}
		plan.textHTML = html

	case ".md":
		html, locals, err := processMd(srcPath)
		if err != nil {
			return postPlan{}, err
		}
		plan.mdHTML = html
		plan.localImages = locals

	case ".png", ".jpg", ".jpeg", ".gif":
		img, err := validateImage(srcPath)
		if err != nil {
			return postPlan{}, err
		}
		plan.validatedImage = img

	case ".mp3":
		if err := validateAudio(srcPath); err != nil {
			return postPlan{}, err
		}

	default:
		return postPlan{}, fmt.Errorf("unsupported file type: %s", ext)
	}

	return plan, nil
}

// commitPlan generates a unique ID, creates the post directory, writes assets,
// persists the post metadata, and removes the source file.
func commitPlan(plan postPlan, postsDir string, now time.Time) error {
	id, err := uniqueID(postsDir, now)
	if err != nil {
		return fmt.Errorf("generate unique ID: %w", err)
	}

	postDir := filepath.Join(postsDir, id)
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		return fmt.Errorf("create post dir %s: %w", id, err)
	}

	var p *post.Post
	var inboxExtras []string

	switch plan.ext {
	case ".txt":
		p = &post.Post{
			ID:        id,
			Timestamp: now,
			PostType:  post.TypeText,
			Content:   plan.textHTML,
		}

	case ".md":
		html := plan.mdHTML
		for _, name := range plan.localImages {
			html = strings.ReplaceAll(html,
				fmt.Sprintf(`src="%s"`, name),
				fmt.Sprintf(`src="posts/%s/%s"`, id, name))
		}

		sourceDir := filepath.Dir(plan.srcPath)
		for _, name := range plan.localImages {
			src := filepath.Join(sourceDir, name)
			dst := filepath.Join(postDir, name)
			if err := copyFile(src, dst); err != nil {
				_ = os.RemoveAll(postDir)
				return fmt.Errorf("copy markdown asset %s: %w", name, err)
			}
			inboxExtras = append(inboxExtras, src)
		}

		p = &post.Post{
			ID:        id,
			Timestamp: now,
			PostType:  post.TypeMarkdown,
			Content:   html,
			Assets:    plan.localImages,
		}

	case ".png", ".jpg", ".jpeg", ".gif":
		if err := writeImageAsset(plan.validatedImage, postDir); err != nil {
			_ = os.RemoveAll(postDir)
			return err
		}
		src := fmt.Sprintf("posts/%s/image.jpg", id)
		html := fmt.Sprintf(`<img src="%s" alt="" class="post-image">`, src)
		p = &post.Post{
			ID:        id,
			Timestamp: now,
			PostType:  post.TypeImage,
			Content:   html,
			Assets:    []string{"image.jpg"},
		}

	case ".mp3":
		outName := filepath.Base(plan.srcPath)
		dst := filepath.Join(postDir, outName)
		if err := copyFile(plan.srcPath, dst); err != nil {
			_ = os.RemoveAll(postDir)
			return err
		}
		src := fmt.Sprintf("posts/%s/%s", id, outName)
		html := fmt.Sprintf(
			`<audio controls class="post-audio"><source src="%s" type="audio/mpeg">Your browser does not support audio.</audio>`,
			src,
		)
		p = &post.Post{
			ID:        id,
			Timestamp: now,
			PostType:  post.TypeAudio,
			Content:   html,
			Assets:    []string{outName},
		}
	}

	if err := p.Save(postDir); err != nil {
		_ = os.RemoveAll(postDir)
		return err
	}

	for _, path := range inboxExtras {
		_ = os.Remove(path)
	}

	return os.Remove(plan.srcPath)
}

// claimedByMarkdown scans all .md entries in inputDir and returns a set of
// image filenames that are referenced within those markdown files.
// Those images should be embedded in the markdown post, not processed alone.
// If two different markdown files claim the same image, an error is returned.
func claimedByMarkdown(entries []os.DirEntry, inputDir string) (map[string]bool, error) {
	claimed := make(map[string]bool)
	// owners tracks which markdown file first claimed each image so we can
	// detect conflicts before processing begins.
	owners := make(map[string]string)

	for _, entry := range entries {
		if entry.IsDir() || strings.ToLower(filepath.Ext(entry.Name())) != ".md" {
			continue
		}

		mdPath := filepath.Join(inputDir, entry.Name())
		data, err := os.ReadFile(mdPath)
		if err != nil {
			return nil, fmt.Errorf("read markdown for image claims %s: %w", entry.Name(), err)
		}

		for _, imgName := range findLocalImages(string(data), inputDir) {
			if owner, exists := owners[imgName]; exists && owner != entry.Name() {
				return nil, fmt.Errorf("image %q claimed by both %q and %q", imgName, owner, entry.Name())
			}
			owners[imgName] = entry.Name()
			claimed[imgName] = true
		}
	}

	return claimed, nil
}

// uniqueID generates a post ID for the given time that does not already exist
// as a directory under postsDir. Appends a numeric suffix if needed.
func uniqueID(postsDir string, t time.Time) (string, error) {
	for i := 0; ; i++ {
		id := post.NewID(t, i)
		_, err := os.Stat(filepath.Join(postsDir, id))
		if err != nil {
			if os.IsNotExist(err) {
				return id, nil
			}
			return "", fmt.Errorf("stat post dir %s: %w", id, err)
		}
	}
}
