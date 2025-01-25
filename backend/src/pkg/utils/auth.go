// ./src/pkg/utils/auth.go
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

type TokenManager struct {
	secretKey    []byte
	rotationTime time.Duration
	timeProvider func() time.Time
}

// NewTokenManager creates a new token manager with the given secret key and rotation interval
func NewTokenManager(secretKey string, rotationTime time.Duration) *TokenManager {
	return &TokenManager{
		secretKey:    []byte(secretKey),
		rotationTime: rotationTime,
		timeProvider: time.Now,
	}
}

// GetRotationTime returns the token rotation time
func (tm *TokenManager) GetRotationTime() time.Duration {
	return tm.rotationTime
}

// GenerateToken creates a new token based on the current time interval
func (tm *TokenManager) GenerateToken() string {
	currentInterval := tm.GetCurrentInterval()
	return tm.generateTokenForInterval(currentInterval)
}

// ValidateToken checks if a token is valid for either the current or previous time interval
func (tm *TokenManager) ValidateToken(token string) bool {
	currentInterval := tm.GetCurrentInterval()

	// Check current interval
	if token == tm.generateTokenForInterval(currentInterval) {
		return true
	}

	// Check previous interval to allow for clock skew
	if token == tm.generateTokenForInterval(currentInterval-1) {
		return true
	}

	return false
}

// getCurrentInterval returns the current time interval number
func (tm *TokenManager) GetCurrentInterval() int64 {
	return tm.timeProvider().Unix() / int64(tm.rotationTime.Seconds())
}

// generateTokenForInterval creates a token for a specific time interval
func (tm *TokenManager) generateTokenForInterval(interval int64) string {
	h := hmac.New(sha256.New, tm.secretKey)
	h.Write([]byte(fmt.Sprintf("%d", interval)))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
