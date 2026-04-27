package processor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/snonux/snonux/internal/post"
)

type audioBuilder struct{}

func (audioBuilder) Plan(srcPath string, ext string) (postPlan, error) {
	plan := postPlan{srcPath: srcPath, ext: ext}
	if err := validateAudio(srcPath); err != nil {
		return postPlan{}, err
	}
	return plan, nil
}

func (audioBuilder) Commit(plan postPlan, postDir string, id string, now time.Time) (*post.Post, []string, error) {
	outName := filepath.Base(plan.srcPath)
	dst := filepath.Join(postDir, outName)
	if err := copyFile(plan.srcPath, dst); err != nil {
		return nil, nil, err
	}
	src := fmt.Sprintf("posts/%s/%s", id, outName)
	html := fmt.Sprintf(
		`<audio controls class="post-audio"><source src="%s" type="audio/mpeg">Your browser does not support audio.</audio>`,
		src,
	)
	p := &post.Post{
		ID:        id,
		Timestamp: now,
		PostType:  post.TypeAudio,
		Content:   html,
		Assets:    []string{outName},
	}
	return p, nil, nil
}

func init() {
	register(".mp3", audioBuilder{})
}

// validateAudio confirms the audio source file exists and is readable.
func validateAudio(srcPath string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open audio %s: %w", srcPath, err)
	}
	return f.Close()
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
