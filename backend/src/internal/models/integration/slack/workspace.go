package slack

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/wizenheimer/byrd/src/internal/config"
)

// SlackWorkspaceStatus represents the status of the slack workspace
type SlackWorkspaceStatus string

const (
	// SlackWorkspaceStatusActive represents the active status of the slack workspace
	// This is the status when workspace is linked and channel and canvas are set
	// This is the final outcome of linking a workspace
	SlackWorkspaceStatusActive SlackWorkspaceStatus = "active"

	// SlackWorkspaceStatusInactive represents the inactive status of the slack workspace
	// This is the status when workspace is unlinked
	// This is the outcome of unlinking a workspace
	SlackWorkspaceStatusInactive SlackWorkspaceStatus = "inactive"
)

// SlackWorkspace represents a slack workspace
type SlackWorkspace struct {
	// Reference to an existing workspace ID
	WorkspaceID uuid.UUID `json:"workspace_id"`

	// Reference to slack team ID
	TeamID string `json:"team_id"`

	// Reference to slack channel ID
	ChannelID string `json:"channel_id"`

	// Reference to slack channel webhook URL
	ChannelWebhookURL string `json:"channel_webhook"`

	// Slack Access Token
	AccessToken string `json:"access_token"`

	// Status of the workspace
	Status SlackWorkspaceStatus `json:"status"`

	// CreatedAt is the time the page was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is the time the page was last updated
	UpdatedAt time.Time `json:"updated_at"`

	// Internal field to track if sensitive data is decoded
	// Not stored in database
	IsDecoded bool `json:"-"`
}

// Encode encodes the sensitive fields before saving to the database
func (s *SlackWorkspace) Encode() error {
	// If already encoded, skip
	if !s.IsDecoded {
		return nil
	}

	if s.AccessToken == "" {
		return fmt.Errorf("access token isn't set")
	}

	// Encrypt access token
	encryptedToken, err := encrypt(s.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}
	s.AccessToken = encryptedToken
	s.IsDecoded = false

	return nil
}

// Decode decodes the sensitive fields before returning to the user
func (s *SlackWorkspace) Decode() error {
	// If already decoded, skip
	if s.IsDecoded {
		return nil
	}

	if s.AccessToken == "" {
		return fmt.Errorf("access token isn't set")
	}

	// Decrypt access token
	decryptedToken, err := decrypt(s.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt access token: %w", err)
	}
	s.AccessToken = decryptedToken
	s.IsDecoded = true

	return nil
}

// encrypt encrypts data using AES-256-GCM with a random nonce
func encrypt(data string) (string, error) {
	// Convert string to bytes
	plaintext := []byte(data)

	// Get secret key
	secretKeyString, err := config.GetSecretKey()
	if err != nil {
		return "", err
	}

	secretKey := []byte(secretKeyString)

	// Create cipher block
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and seal data
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// Encode to base64 for storage
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return encoded, nil
}

// decrypt decrypts data using AES-256-GCM
func decrypt(encodedData string) (string, error) {
	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Get Secret Key
	secretKeyString, err := config.GetSecretKey()
	if err != nil {
		return "", err
	}

	secretKey := []byte(secretKeyString)

	// Create cipher block
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce from ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
