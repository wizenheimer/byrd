package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/slack-go/slack"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type SlackIntegrationHandler struct {
	store  *session.Store
	logger *logger.Logger
}

func NewSlackIntegrationHandler(
	logger *logger.Logger,
	redisURL string,
) (*SlackIntegrationHandler, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	if redisURL == "" {
		return nil, errors.New("redisURL is required")
	}

	store := redis.New(redis.Config{
		URL: redisURL,
	})

	sessionStore := session.New(session.Config{
		Storage:    store,
		Expiration: 6000, // 10 minutes
	})

	h := SlackIntegrationHandler{
		store: sessionStore,
		logger: logger.WithFields(
			map[string]any{
				"module": "slack_integration_handler",
			},
		),
	}

	return &h, nil
}

func (sh *SlackIntegrationHandler) SlackInstallationHandler(c *fiber.Ctx) error {
	sess, err := sh.store.Get(c)
	if err != nil {
		return c.Status(500).SendString("Failed to get session")
	}

	code := c.Query("code")
	stateToken := c.Query("state")

	if code == "" || stateToken == "" {
		return c.Status(400).SendString("Missing parameters")
	}

	// Retrieve state from session
	encodedState := sess.Get(stateToken)
	if encodedState == nil {
		return c.Status(400).SendString("Invalid or expired state")
	}

	// Decode state from Base64
	stateBytes, err := base64.URLEncoding.DecodeString(encodedState.(string))
	if err != nil {
		return c.Status(400).SendString("Invalid state format")
	}

	var oauthState OAuthState
	if err := json.Unmarshal(stateBytes, &oauthState); err != nil {
		return c.Status(400).SendString("Invalid state format")
	}

	resp, err := slack.GetOAuthV2Response(http.DefaultClient,
		os.Getenv("SLACK_CLIENT_ID"),
		os.Getenv("SLACK_CLIENT_SECRET"),
		code,
		os.Getenv("SLACK_REDIRECT_URI"),
	)

	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Error getting Slack token: %s", err.Error()))
	}

	// Store the token in the database
	fmt.Println("Slack token:", resp.AccessToken, "for workspace:", oauthState.WorkspaceID)

	// Clear state after use
	sess.Delete(stateToken)
	if err := sess.Save(); err != nil {
		sh.logger.Error("failed to save session post deletion", zap.Error(err))
	}

	// Trigger a redirect
	return c.Redirect(fmt.Sprintf("slack://app?team=%s&id=%s&tab=about", resp.Team.ID, resp.AppID), http.StatusFound)
}

func (sh *SlackIntegrationHandler) SlackUninstallationHandler(c *fiber.Ctx) error {
	return nil
}

func (sh *SlackIntegrationHandler) StatusSlackHandler(c *fiber.Ctx) error {
	return nil
}

func (sh *SlackIntegrationHandler) SlackOAuthHandler(c *fiber.Ctx) error {
	sess, err := sh.store.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get session"})
	}

	var req struct {
		WorkspaceID string `json:"workspace_id"`
		UserID      string `json:"user_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.WorkspaceID == "" || req.UserID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing parameters"})
	}

	// Generate a secure random token
	token, err := generateToken()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Create a state struct
	oauthState := OAuthState{
		WorkspaceID: req.WorkspaceID,
		UserID:      req.UserID,
		Token:       token,
	}

	// Encode state as JSON and Base64
	stateBytes, err := json.Marshal(oauthState)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to encode state"})
	}

	encodedState := base64.URLEncoding.EncodeToString(stateBytes)

	// Store state in session
	sess.Set(oauthState.Token, encodedState)
	if err := sess.Save(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save session"})
	}

	slackOAuthURL := fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&scope=canvases:read,canvases:write,channels:read,chat:write,chat:write.public,commands,im:write,users.profile:read&state=%s&redirect_uri=%s",
		os.Getenv("SLACK_CLIENT_ID"),
		oauthState.Token,
		os.Getenv("SLACK_REDIRECT_URI"),
	)

	return sendDataResponse(
		c,
		http.StatusOK,
		"Slack OAuth URL",
		map[string]string{
			"url": slackOAuthURL,
		},
	)
}

// State struct to encode workspace & user info
type OAuthState struct {
	WorkspaceID string `json:"workspace_id"`
	UserID      string `json:"user_id"`
	Token       string `json:"token"`
}

// Generate a secure random token
func generateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
