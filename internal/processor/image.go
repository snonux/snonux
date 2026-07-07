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
// It sniffs the actual file type from magic bytes instead of trusting the extension.
func decodeImage(srcPath string) (image.Image, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("open image %s: %w", srcPath, err)
	}
	defer f.Close()

	// Read the first 512 bytes to detect the actual file format.
	const sniffLen = 512
	head := make([]byte, sniffLen)
	n, err := f.Read(head)
	if err != nil {
		return nil, fmt.Errorf("read image %s: %w", srcPath, err)
	}
	if n == 0 {
		return nil, fmt.Errorf("empty image file: %s", srcPath)
	}

	// Determine the actual image format from magic bytes.
	sniff := head[:n]
	var format string
	switch {
	case len(sniff) >= 2 && sniff[0] == 0xFF && sniff[1] == 0xD8:
		format = "jpeg"
	case len(sniff) >= 8 && string(sniff[:8]) == "\x89PNG\r\n\x1a\n":
		format = "png"
	case len(sniff) >= 6 && string(sniff[:6]) == "GIF87a" || string(sniff[:6]) == "GIF89a":
		format = "gif"
	default:
		return nil, fmt.Errorf("unsupported image format for %s (extension: %s sniff: %q)", srcPath, filepath.Ext(srcPath), string(sniff))
	}

	// Rewind to the beginning so the decoder sees the full stream.
	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("rewind image %s: %w", srcPath, err)
	}

	var img image.Image
	var decodeErr error
	switch format {
	case "jpeg":
		img, decodeErr = jpeg.Decode(f)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode JPEG %s: %w", srcPath, decodeErr)
		}
	case "png":
		img, decodeErr = png.Decode(f)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode PNG %s: %w", srcPath, decodeErr)
		}
	case "gif":
		// Use only the first frame of animated GIFs.
		img, decodeErr = gif.Decode(f)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode GIF %s: %w", srcPath, decodeErr)
		}
	}

	return img, nil
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
