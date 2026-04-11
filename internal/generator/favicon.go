package generator

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strings"
)

const faviconHeadHTML = `
    <link rel="icon" href="favicon.ico" sizes="any">
`

func injectSharedHead(theme string) string {
	if strings.Contains(theme, `rel="icon"`) {
		return theme
	}

	return strings.Replace(theme, "</head>", faviconHeadHTML+"</head>", 1)
}

func writeFavicon(outputDir string) error {
	data, err := generateFaviconICO()
	if err != nil {
		return fmt.Errorf("generate favicon.ico: %w", err)
	}

	if err := os.WriteFile(filepath.Join(outputDir, "favicon.ico"), data, 0o644); err != nil {
		return fmt.Errorf("write favicon.ico: %w", err)
	}

	return nil
}

func generateFaviconICO() ([]byte, error) {
	const size = 32

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	black := color.RGBA{R: 8, G: 10, B: 12, A: 255}
	green := color.RGBA{R: 0, G: 255, B: 102, A: 255}
	cyan := color.RGBA{R: 0, G: 214, B: 255, A: 255}

	fillRect(img, 0, 0, size, size, black)
	fillRect(img, 1, 1, size-1, 2, green)
	fillRect(img, 1, size-2, size-1, size-1, green)
	fillRect(img, 1, 1, 2, size-1, green)
	fillRect(img, size-2, 1, size-1, size-1, cyan)

	// Left glyph: blocky "S".
	fillRect(img, 5, 5, 13, 8, green)
	fillRect(img, 5, 8, 8, 14, green)
	fillRect(img, 5, 14, 13, 17, green)
	fillRect(img, 10, 17, 13, 23, green)
	fillRect(img, 5, 23, 13, 26, green)

	// Right glyph: blocky "N".
	fillRect(img, 18, 5, 21, 26, cyan)
	fillRect(img, 25, 5, 28, 26, cyan)
	for i := 0; i < 7; i++ {
		x := 20 + i
		y := 6 + i*3
		fillRect(img, x, y, x+2, y+3, cyan)
	}

	return encodeICO(img)
}

func fillRect(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	bounds := img.Bounds()
	if x0 < bounds.Min.X {
		x0 = bounds.Min.X
	}
	if y0 < bounds.Min.Y {
		y0 = bounds.Min.Y
	}
	if x1 > bounds.Max.X {
		x1 = bounds.Max.X
	}
	if y1 > bounds.Max.Y {
		y1 = bounds.Max.Y
	}

	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.SetRGBA(x, y, c)
		}
	}
}

func encodeICO(img image.Image) ([]byte, error) {
	b := img.Bounds()
	if b.Dx() != 32 || b.Dy() != 32 {
		return nil, fmt.Errorf("favicon must be 32x32, got %dx%d", b.Dx(), b.Dy())
	}

	const (
		bitmapHeaderSize = 40
		icoHeaderSize    = 6
		dirEntrySize     = 16
		bitsPerPixel     = 32
	)

	maskRowSize := ((b.Dx() + 31) / 32) * 4
	maskSize := maskRowSize * b.Dy()
	pixelSize := b.Dx() * b.Dy() * 4
	imageSize := bitmapHeaderSize + pixelSize + maskSize
	imageOffset := icoHeaderSize + dirEntrySize

	var buf bytes.Buffer

	write := func(v any) error {
		return binary.Write(&buf, binary.LittleEndian, v)
	}

	if err := write(uint16(0)); err != nil {
		return nil, err
	}
	if err := write(uint16(1)); err != nil {
		return nil, err
	}
	if err := write(uint16(1)); err != nil {
		return nil, err
	}

	if err := buf.WriteByte(byte(b.Dx())); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(byte(b.Dy())); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(0); err != nil {
		return nil, err
	}
	if err := buf.WriteByte(0); err != nil {
		return nil, err
	}
	if err := write(uint16(1)); err != nil {
		return nil, err
	}
	if err := write(uint16(bitsPerPixel)); err != nil {
		return nil, err
	}
	if err := write(uint32(imageSize)); err != nil {
		return nil, err
	}
	if err := write(uint32(imageOffset)); err != nil {
		return nil, err
	}

	if err := write(uint32(bitmapHeaderSize)); err != nil {
		return nil, err
	}
	if err := write(int32(b.Dx())); err != nil {
		return nil, err
	}
	if err := write(int32(b.Dy() * 2)); err != nil {
		return nil, err
	}
	if err := write(uint16(1)); err != nil {
		return nil, err
	}
	if err := write(uint16(bitsPerPixel)); err != nil {
		return nil, err
	}
	if err := write(uint32(0)); err != nil {
		return nil, err
	}
	if err := write(uint32(pixelSize + maskSize)); err != nil {
		return nil, err
	}
	if err := write(int32(0)); err != nil {
		return nil, err
	}
	if err := write(int32(0)); err != nil {
		return nil, err
	}
	if err := write(uint32(0)); err != nil {
		return nil, err
	}
	if err := write(uint32(0)); err != nil {
		return nil, err
	}

	for y := b.Max.Y - 1; y >= b.Min.Y; y-- {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, a := img.At(x, y).RGBA()
			if err := buf.WriteByte(byte(bl >> 8)); err != nil {
				return nil, err
			}
			if err := buf.WriteByte(byte(g >> 8)); err != nil {
				return nil, err
			}
			if err := buf.WriteByte(byte(r >> 8)); err != nil {
				return nil, err
			}
			if err := buf.WriteByte(byte(a >> 8)); err != nil {
				return nil, err
			}
		}
	}

	if _, err := buf.Write(make([]byte, maskSize)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
