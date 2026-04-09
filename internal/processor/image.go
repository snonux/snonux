package processor

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

const (
	maxImageWidth  = 1024
	jpegQuality    = 80
)

// processImage reads the source image, resizes it if wider than maxImageWidth,
// encodes it as JPEG at jpegQuality, and writes the result to destDir.
// Returns the output filename (always a .jpg) and an HTML <img> snippet.
func processImage(srcPath, destDir, postID string) (filename, htmlContent string, err error) {
	img, err := decodeImage(srcPath)
	if err != nil {
		return "", "", err
	}

	img = resizeIfNeeded(img)

	outName := "image.jpg"
	outPath := filepath.Join(destDir, outName)

	if err := writeJPEG(img, outPath); err != nil {
		return "", "", err
	}

	// The <img> src is relative to the site root, pointing into the posts dir.
	src := fmt.Sprintf("posts/%s/%s", postID, outName)
	html := fmt.Sprintf(`<img src="%s" alt="" class="post-image">`, src)

	return outName, html, nil
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
