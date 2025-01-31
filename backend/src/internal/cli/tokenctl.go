package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/wizenheimer/byrd/src/pkg/utils"
)

type TokenResponse struct {
	Token        string        `json:"token"`
	ExpiresIn    int64         `json:"expires_in"`
	RefreshAfter time.Duration `json:"refresh_after"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var tokenManager *utils.TokenManager

func main() {
	// Load environment variables
	if err := godotenv.Load(".env.development"); err != nil {
		log.Fatal("Error loading .env.development file")
	}

	secretKey := os.Getenv("MANAGEMENT_API_KEY")
	if secretKey == "" {
		log.Fatal("MANAGEMENT_API_KEY is required")
	}

	// Initialize token manager
	rotationTime := 10 * time.Minute // Default rotation time
	tokenManager = utils.NewTokenManager(secretKey, rotationTime)

	// Create new Fiber app
	app := fiber.New()

	// Set up routes
	app.Get("/token", handleGetToken)
	app.Get("/token/status", handleTokenStatus)
	app.Get("/health", handleHealth)

	// Start server
	port := 4004
	fmt.Printf("Starting token management server on port %d...\n", port)
	if err := app.Listen(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
}

func handleGetToken(c *fiber.Ctx) error {
	token := tokenManager.GenerateToken()
	remaining := getRemainingTime()

	return c.JSON(TokenResponse{
		Token:        token,
		ExpiresIn:    int64(remaining.Seconds()),
		RefreshAfter: remaining,
	})
}

func handleTokenStatus(c *fiber.Ctx) error {
	remaining := getRemainingTime()
	currentInterval := tokenManager.GetCurrentInterval()

	return c.JSON(fiber.Map{
		"current_interval": currentInterval,
		"expires_in":       remaining.Seconds(),
		"rotation_time":    tokenManager.GetRotationTime(),
	})
}

func handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
	})
}

func getRemainingTime() time.Duration {
	currentInterval := tokenManager.GetCurrentInterval()
	rotationTime := tokenManager.GetRotationTime()
	nextRotation := (currentInterval + 1) * int64(rotationTime.Seconds())
	return time.Duration(nextRotation-time.Now().Unix()) * time.Second
}
