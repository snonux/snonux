package processor

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"

	"codeberg.org/snonux/snonux/internal/post"
	"golang.org/x/image/draw"
)

const (
	maxImageWidth = 1024
	jpegQuality   = 80
)

type imageBuilder struct{}

func (imageBuilder) Plan(srcPath string, ext string) (postPlan, error) {
	plan := postPlan{srcPath: srcPath, ext: ext}
	img, err := validateImage(srcPath)
	if err != nil {
		return postPlan{}, err
	}
	plan.validatedImage = img
	return plan, nil
}

func (imageBuilder) Commit(plan postPlan, postDir string, id string, now time.Time) (*post.Post, []string, error) {
	if err := writeImageAsset(plan.validatedImage, postDir); err != nil {
		return nil, nil, err
	}
	src := fmt.Sprintf("posts/%s/image.jpg", id)
	html := fmt.Sprintf(`<img src="%s" alt="" class="post-image">`, src)
	p := &post.Post{
		ID:        id,
		Timestamp: now,
		PostType:  post.TypeImage,
		Content:   html,
		Assets:    []string{"image.jpg"},
	}
	return p, nil, nil
}

func init() {
	register(".png", imageBuilder{})
	register(".jpg", imageBuilder{})
	register(".jpeg", imageBuilder{})
	register(".gif", imageBuilder{})
}

// validateImage reads and decodes the source image, resizing if necessary.
// It performs only read validation; the caller is responsible for writing assets.
func validateImage(srcPath string) (image.Image, error) {
	img, err := decodeImage(srcPath)
	if err != nil {
		return nil, err
	}
	return resizeIfNeeded(img), nil
}

// writeImageAsset writes the prepared image as JPEG into postDir.
func writeImageAsset(img image.Image, postDir string) error {
	outPath := filepath.Join(postDir, "image.jpg")
	return writeJPEG(img, outPath)
}

// decodeImage decodes a JPEG, PNG, or GIF (first frame) from srcPath.
func decodeImage(srcPath string) (image.Image, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("open image %s: %w", srcPath, err)
	}
	defer f.Close()

	ext := filepath.Ext(srcPath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err := jpeg.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("decode JPEG %s: %w", srcPath, err)
		}
		return img, nil

	case ".png":
		img, err := png.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("decode PNG %s: %w", srcPath, err)
		}
		return img, nil

	case ".gif":
		// Use only the first frame of animated GIFs.
		g, err := gif.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("decode GIF %s: %w", srcPath, err)
		}
		return g, nil

	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
}

// resizeIfNeeded returns a resized copy of img if its width exceeds maxImageWidth,
// preserving aspect ratio. Otherwise the original is returned unchanged.
func resizeIfNeeded(img image.Image) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()

	if w <= maxImageWidth {
		return img
	}

	h := bounds.Dy()
	newW := maxImageWidth
	newH := (h * newW) / w

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	return dst
}

// writeJPEG encodes img as JPEG at the configured quality level and writes to path.
func writeJPEG(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create JPEG %s: %w", path, err)
	}
	defer f.Close()

	opts := &jpeg.Options{Quality: jpegQuality}
	if err := jpeg.Encode(f, img, opts); err != nil {
		return fmt.Errorf("encode JPEG %s: %w", path, err)
	}

	return nil
}
