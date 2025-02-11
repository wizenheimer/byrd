package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type SlackIntegrationHandler struct {
	logger *logger.Logger
}

func NewSlackIntegrationHandler(
	logger *logger.Logger,
) (*SlackIntegrationHandler, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	h := SlackIntegrationHandler{
		logger: logger.WithFields(
			map[string]any{
				"module": "slack_integration_handler",
			},
		),
	}

	return &h, nil
}

// generateToken generates a secure random token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// SlackOAuthHandler initiates the Slack OAuth flow
func (sh *SlackIntegrationHandler) SlackOAuthHandler(c *fiber.Ctx) error {
	// Validate required parameters
	workspaceID := c.Query("workspace_id", "")
	if workspaceID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing workspace_id"})
	}

	userID := c.Query("user_id", "")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing user_id"})
	}

	// Generate token
	token, err := generateToken()
	if err != nil {
		sh.logger.Error("Failed to generate token", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Create state object
	oauthState := OAuthState{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Token:       token,
		CreatedAt:   time.Now().UTC(),
	}

	// Encode state
	stateBytes, err := json.Marshal(oauthState)
	if err != nil {
		sh.logger.Error("Failed to marshal state", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": "Failed to encode state"})
	}
	encodedState := base64.URLEncoding.EncodeToString(stateBytes)

	// Set cookie with state
	cookie := fiber.Cookie{
		Name:     "slack_oauth_state",
		Value:    encodedState,
		Path:     "/",
		MaxAge:   900, // 15 minutes
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	}
	c.Cookie(&cookie)

	// Generate Slack OAuth URL
	scopes := os.Getenv("SLACK_CLIENT_SCOPES")
	if scopes == "" {
		scopes = "channels:read,channels:join,chat:write,commands" // default scopes
	}

	slackOAuthURL := fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&scope=%s&state=%s&redirect_uri=%s",
		os.Getenv("SLACK_CLIENT_ID"),
		scopes,
		encodedState,
		os.Getenv("SLACK_REDIRECT_URI"),
	)

	return c.Redirect(slackOAuthURL, http.StatusFound)
}

// SlackInstallationHandler handles the OAuth callback from Slack
func (sh *SlackIntegrationHandler) SlackInstallationHandler(c *fiber.Ctx) error {
	// Get OAuth parameters
	code := c.Query("code")
	stateToken := c.Query("state")

	if code == "" || stateToken == "" {
		sh.logger.Error("Missing OAuth parameters")
		return c.Status(400).SendString("Missing parameters")
	}

	// Get state from cookie
	stateCookie := c.Cookies("slack_oauth_state")
	if stateCookie == "" {
		sh.logger.Error("State cookie not found",
			zap.String("receivedState", stateToken))
		return c.Status(400).SendString("Invalid or expired state")
	}

	// Compare received state with cookie state
	if stateCookie != stateToken {
		sh.logger.Error("State mismatch",
			zap.String("receivedState", stateToken),
			zap.String("cookieState", stateCookie))
		return c.Status(400).SendString("Invalid state")
	}

	// Decode state
	stateBytes, err := base64.URLEncoding.DecodeString(stateCookie)
	if err != nil {
		sh.logger.Error("Failed to decode state", zap.Error(err))
		return c.Status(400).SendString("Invalid state format")
	}

	var oauthState OAuthState
	if err := json.Unmarshal(stateBytes, &oauthState); err != nil {
		sh.logger.Error("Failed to unmarshal state", zap.Error(err))
		return c.Status(400).SendString("Invalid state format")
	}

	// Verify state not expired (15 minutes)
	if time.Since(oauthState.CreatedAt) > 15*time.Minute {
		sh.logger.Error("State expired",
			zap.Time("createdAt", oauthState.CreatedAt),
			zap.Duration("age", time.Since(oauthState.CreatedAt)))
		return c.Status(400).SendString("State expired")
	}

	// Exchange code for token
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := slack.GetOAuthV2Response(
		client,
		os.Getenv("SLACK_CLIENT_ID"),
		os.Getenv("SLACK_CLIENT_SECRET"),
		code,
		os.Getenv("SLACK_REDIRECT_URI"),
	)
	if err != nil {
		sh.logger.Error("Slack OAuth exchange failed", zap.Error(err))
		return c.Status(500).SendString(fmt.Sprintf("Error getting Slack token: %s", err.Error()))
	}

	// Store integration details
	integration := SlackIntegration{
		WorkspaceID:      oauthState.WorkspaceID,
		UserID:           oauthState.UserID,
		SlackAccessToken: resp.AccessToken,
		SlackTeamID:      resp.Team.ID,
		SlackAppID:       resp.AppID,
		CreatedAt:        time.Now().UTC(),
	}

	// TODO: Store integration in your database
	sh.logger.Info("Slack integration successful",
		zap.String("workspaceID", integration.WorkspaceID),
		zap.String("userID", integration.UserID),
		zap.String("teamID", integration.SlackTeamID),
		zap.String("appID", integration.SlackAppID),
		zap.Time("createdAt", integration.CreatedAt),
		zap.String("accessToken", integration.SlackAccessToken),
	)

	// Clear the state cookie
	c.Cookie(&fiber.Cookie{
		Name:     "slack_oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	// Redirect back to Slack
	return c.Redirect(
		fmt.Sprintf("slack://app?team=%s&id=%s&tab=about",
			resp.Team.ID,
			resp.AppID,
		),
		http.StatusFound,
	)
}

// SlackIntegration represents the stored integration details
type SlackIntegration struct {
	WorkspaceID      string    `json:"workspace_id"`
	UserID           string    `json:"user_id"`
	SlackAccessToken string    `json:"slack_access_token"`
	SlackTeamID      string    `json:"slack_team_id"`
	SlackAppID       string    `json:"slack_app_id"`
	CreatedAt        time.Time `json:"created_at"`
}

// OAuthState represents the state object used in OAuth flow
type OAuthState struct {
	WorkspaceID string    `json:"workspace_id"`
	UserID      string    `json:"user_id"`
	Token       string    `json:"token"`
	CreatedAt   time.Time `json:"created_at"`
}
