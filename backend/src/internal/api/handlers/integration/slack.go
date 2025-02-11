package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
	"github.com/valyala/fasthttp"
	slackworkspace "github.com/wizenheimer/byrd/src/internal/service/integration/slack"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type SlackIntegrationHandler struct {
	logger *logger.Logger
	svc    slackworkspace.SlackWorkspaceService
}

func NewSlackIntegrationHandler(
	logger *logger.Logger,
	svc slackworkspace.SlackWorkspaceService,
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
		svc: svc,
	}

	return &h, nil
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
	oauthState := SlackOAuthState{
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

	var oauthState SlackOAuthState
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
	// integration := SlackIntegration{
	// 	WorkspaceID:      oauthState.WorkspaceID,
	// 	UserID:           oauthState.UserID,
	// 	SlackAccessToken: resp.AccessToken,
	// 	SlackTeamID:      resp.Team.ID,
	// 	SlackAppID:       resp.AppID,
	// 	CreatedAt:        time.Now().UTC(),
	// }

	// TODO: Store integration in your database
	workspaceUUID, err := uuid.Parse(oauthState.WorkspaceID)
	if err != nil {
		return c.Status(500).SendString("Failed to parse workspace ID")
	}

	_, err = sh.svc.CreateSlackWorkspace(c.Context(), workspaceUUID, resp.Team.ID, resp.AccessToken)
	if err != nil {
		sh.logger.Error("Failed to create Slack workspace", zap.Error(err))
		return c.Status(500).SendString("Failed to create Slack workspace")
	}

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

// SlackConfigurationHandler handles the configuration of the Slack app
func (sh *SlackIntegrationHandler) ConfigureCommandHandler(c *fiber.Ctx) error {
	cmd, err := SlashCommandParseFast(c.Request())
	if err != nil {
		return c.Status(400).SendString("Failed to parse command")
	}

	_, err = sh.svc.UpdateSlackWorkspace(c.Context(), cmd)
	if err != nil {
		sh.logger.Error("Failed to update Slack workspace", zap.Error(err))
		return c.Status(500).SendString("Failed to update Slack workspace")
	}

	return c.Status(200).Send(nil)
}

func (sh *SlackIntegrationHandler) WatchCommandHandler(c *fiber.Ctx) error {
	cmd, err := SlashCommandParseFast(c.Request())
	if err != nil {
		return c.Status(400).SendString("Failed to parse command")
	}

	if err := sh.svc.CreateCompetitorForWorkspace(c.Context(), cmd); err != nil {
		sh.logger.Error("Failed to create competitor for workspace", zap.Error(err))
		return c.Status(500).SendString("Failed to create competitor for workspace")
	}

	return c.Status(200).Send(nil)
}

func (sh *SlackIntegrationHandler) UserCommandHandler(c *fiber.Ctx) error {
	cmd, err := SlashCommandParseFast(c.Request())
	if err != nil {
		return c.Status(400).SendString("Failed to parse command")
	}

	if err := sh.svc.AddUserToSlackWorkspace(c.Context(), cmd); err != nil {
		sh.logger.Error("Failed to add user to Slack workspace", zap.Error(err))
		return c.Status(500).SendString("Failed to add user to Slack workspace")
	}

	return c.Status(200).Send(nil)
}

func (sh *SlackIntegrationHandler) SlackInteractionHandler(c *fiber.Ctx) error {

	// Get payload from form
	payloadStr := c.FormValue("payload")

	var payload slack.InteractionCallback
	if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	if err := sh.svc.HandleSlackInteractionPayload(c.Context(), payload); err != nil {
		sh.logger.Error("Failed to handle interaction", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to handle interaction"})
	}

	return c.Status(200).Send(nil)
}

func SlashCommandParseFast(req *fasthttp.Request) (s slack.SlashCommand, err error) {
	// Get POST form arguments
	args := req.PostArgs()

	s.Token = string(args.Peek("token"))
	s.TeamID = string(args.Peek("team_id"))
	s.TeamDomain = string(args.Peek("team_domain"))
	s.EnterpriseID = string(args.Peek("enterprise_id"))
	s.EnterpriseName = string(args.Peek("enterprise_name"))
	s.IsEnterpriseInstall = string(args.Peek("is_enterprise_install")) == "true"
	s.ChannelID = string(args.Peek("channel_id"))
	s.ChannelName = string(args.Peek("channel_name"))
	s.UserID = string(args.Peek("user_id"))
	s.UserName = string(args.Peek("user_name"))
	s.Command = string(args.Peek("command"))
	s.Text = string(args.Peek("text"))
	s.ResponseURL = string(args.Peek("response_url"))
	s.TriggerID = string(args.Peek("trigger_id"))
	s.APIAppID = string(args.Peek("api_app_id"))

	return s, nil
}
