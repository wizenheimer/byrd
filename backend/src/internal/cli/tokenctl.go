package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// Set up routes
	http.HandleFunc("/token", handleGetToken)
	http.HandleFunc("/token/status", handleTokenStatus)
	http.HandleFunc("/health", handleHealth)

	// Start server
	port := 4004
	fmt.Printf("Starting token management server on port %d...\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

func handleGetToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := tokenManager.GenerateToken()
	remaining := getRemainingTime()

	response := TokenResponse{
		Token:        token,
		ExpiresIn:    int64(remaining.Seconds()),
		RefreshAfter: remaining,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleTokenStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	remaining := getRemainingTime()
	currentInterval := tokenManager.GetCurrentInterval()

	status := map[string]interface{}{
		"current_interval": currentInterval,
		"expires_in":       remaining.Seconds(),
		"rotation_time":    tokenManager.GetRotationTime(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := map[string]string{
		"status": "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func getRemainingTime() time.Duration {
	currentInterval := tokenManager.GetCurrentInterval()
	rotationTime := tokenManager.GetRotationTime()
	nextRotation := (currentInterval + 1) * int64(rotationTime.Seconds())
	return time.Duration(nextRotation-time.Now().Unix()) * time.Second
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
