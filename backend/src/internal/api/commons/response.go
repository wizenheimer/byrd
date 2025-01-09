package commons

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/config"
)

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
	ErrorDetails any    `json:"details,omitempty"` // add omitempty to hide when empty
}

func SendErrorResponse(c *fiber.Ctx, status int, message string, details ...any) error {
	// Check if debug header exists or if the app is in development mode
	isDebug := c.Get("X-Debug") == "true" || config.IsDevelopment()

	response := ErrorResponse{
		ErrorMessage: message,
	}

	// Only include error details if debug header is present
	if isDebug {
		response.ErrorDetails = details
	}

	return c.Status(status).JSON(response)
}

type DataResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SendDataResponse(c *fiber.Ctx, status int, message string, data any) error {
	response := DataResponse{
		Message: message,
	}

	if data != nil {
		response.Data = data
	}

	return c.Status(status).JSON(response)
}
