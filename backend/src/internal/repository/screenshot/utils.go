// ./src/internal/repository/screenshot/utils.go
package screenshot

import (
	"image"
	"io"

	"image/jpeg"
	"image/png"
)

// encodeImage encodes an image.Image to the specified writer
func encodeImage(img image.Image, w io.Writer) error {
	switch v := img.(type) {
	case *image.NRGBA, *image.RGBA:
		return png.Encode(w, v)
	case *image.YCbCr:
		return jpeg.Encode(w, v, &jpeg.Options{Quality: 90})
	default:
		// Default to PNG for unknown types
		return png.Encode(w, img)
	}
}
