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

type TokenServer struct {
	tokenManager *utils.TokenManager
	app          *fiber.App
}

func NewTokenServer(secretKey string, rotationTime time.Duration) (*TokenServer, error) {
	tm := utils.NewTokenManager(secretKey, rotationTime)
	if tm == nil {
		return nil, fmt.Errorf("failed to create token manager")
	}

	app := fiber.New()
	server := &TokenServer{
		tokenManager: tm,
		app:          app,
	}

	// Set up routes
	app.Get("/token", server.handleGetToken)
	app.Get("/token/status", server.handleTokenStatus)
	app.Get("/health", server.handleHealth)

	return server, nil
}

func (s *TokenServer) Start(port int) error {
	fmt.Printf("Starting token management server on port %d...\n", port)
	return s.app.Listen(fmt.Sprintf(":%d", port))
}

func (s *TokenServer) handleGetToken(c *fiber.Ctx) error {
	token := s.tokenManager.GenerateToken()
	remaining := s.getRemainingTime()

	return c.JSON(TokenResponse{
		Token:        token,
		ExpiresIn:    int64(remaining.Seconds()),
		RefreshAfter: remaining,
	})
}

func (s *TokenServer) handleTokenStatus(c *fiber.Ctx) error {
	remaining := s.getRemainingTime()
	currentInterval := s.tokenManager.GetCurrentInterval()

	return c.JSON(fiber.Map{
		"current_interval": currentInterval,
		"expires_in":       remaining.Seconds(),
		"rotation_time":    s.tokenManager.GetRotationTime(),
	})
}

func (s *TokenServer) handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
	})
}

func (s *TokenServer) getRemainingTime() time.Duration {
	currentInterval := s.tokenManager.GetCurrentInterval()
	rotationTime := s.tokenManager.GetRotationTime()
	nextRotation := (currentInterval + 1) * int64(rotationTime.Seconds())
	return time.Duration(nextRotation-time.Now().Unix()) * time.Second
}

func main() {
	// Load environment variables
	if err := godotenv.Load(".env.development"); err != nil {
		log.Fatal("Error loading .env.development file")
	}

	secretKey := os.Getenv("MANAGEMENT_API_KEY")
	if secretKey == "" {
		log.Fatal("MANAGEMENT_API_KEY is required")
	}

	// Create and start token server
	server, err := NewTokenServer(secretKey, 10*time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	// Start server
	if err := server.Start(4004); err != nil {
		log.Fatal(err)
	}
}
