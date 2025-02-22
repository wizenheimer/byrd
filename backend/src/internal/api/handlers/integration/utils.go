package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// generateToken generates a secure random token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// SlackOAuthState represents the state object used in OAuth flow
type SlackOAuthState struct {
	Competitors []string  `json:"competitors"`
	Features    []string  `json:"features"`
	Profiles    []string  `json:"profiles"`
	Token       string    `json:"token"`
	CreatedAt   time.Time `json:"created_at"`
}
