// ./src/pkg/utils/api.go
package utils

import (
	"bytes"
	"image"
	"image/png"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// WritePNGResponse writes an image to a PNG byte array
func WritePNGResponse(img image.Image) ([]byte, error) {
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func QueryIntPtr(c *fiber.Ctx, key string, defaultValue int) *int {
	if value := c.Query(key); value != "" {
		if val, err := strconv.Atoi(value); err == nil {
			return &val
		}
	}
	return &defaultValue
}
