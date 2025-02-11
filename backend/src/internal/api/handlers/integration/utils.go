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
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	Token       string    `json:"token"`
	CreatedAt   time.Time `json:"created_at"`
}
