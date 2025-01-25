// ./src/pkg/utils/api.go
package utils

import (
	"bytes"
	"image"
	"image/png"
)

// WritePNGResponse writes an image to a PNG byte array
func WritePNGResponse(img image.Image) ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
