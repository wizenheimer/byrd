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
	"github.com/slack-go/slack"
	"github.com/valyala/fasthttp"
	api "github.com/wizenheimer/byrd/src/internal/models/api"
	models "github.com/wizenheimer/byrd/src/internal/models/core"
	"github.com/wizenheimer/byrd/src/internal/service/ai"
	slackworkspace "github.com/wizenheimer/byrd/src/internal/service/integration/slack"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
)

type SlackIntegrationHandler struct {
	logger       *logger.Logger
	slackService slackworkspace.SlackWorkspaceService
}

func NewSlackIntegrationHandler(
	logger *logger.Logger,
	slackService slackworkspace.SlackWorkspaceService,
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
		slackService: slackService,
	}

	return &h, nil
}

// SlackOAuthHandler initiates the Slack OAuth flow
func (sh *SlackIntegrationHandler) SlackOAuthHandler(c *fiber.Ctx) error {
	var req api.WorkspaceCreationRequest
	if err := c.BodyParser(&req); err != nil {
		sh.logger.Error("invalid workspace creation request", zap.Error(err))
		return c.Status(400).JSON(fiber.Map{"error": "Invalid workspace creation request"})
	}

	// Generate token
	token, err := generateToken()
	if err != nil {
		sh.logger.Error("failed to generate token", zap.Error(err))
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Create state object
	oauthState := SlackOAuthState{
		Competitors: req.Competitors,
		Features:    req.Features,
		Profiles:    req.Profiles,
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

	// Generate Slack OAuth URL
	scopes := os.Getenv("SLACK_CLIENT_SCOPES")
	if scopes == "" {
		scopes = "canvases:read,canvases:write,channels:join,channels:read,channels:write.topic,chat:write,chat:write.public,commands,groups:read,groups:write.invites,groups:write.topic,im:write,users:read,users:read.email,users.profile:read,incoming-webhook" // default scopes
	}

	slackOAuthURL := fmt.Sprintf(
		"https://slack.com/oauth/v2/authorize?client_id=%s&scope=%s&state=%s&redirect_uri=%s",
		os.Getenv("SLACK_CLIENT_ID"),
		scopes,
		encodedState,
		os.Getenv("SLACK_REDIRECT_URI"),
	)

	return c.Status(200).JSON(fiber.Map{"oauth_url": slackOAuthURL})
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

	// Decode state
	stateBytes, err := base64.URLEncoding.DecodeString(stateToken)
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

	sh.logger.Debug("create workspace using oauth state", zap.Any("state", oauthState))

	oauthState.Profiles, err = ai.Sanitize(oauthState.Profiles)
	if err != nil {
		return err
	}

	var pages []models.PageProps
	for _, competitorURL := range oauthState.Competitors {
		page, err := models.NewPageProps(competitorURL, oauthState.Profiles)
		if err != nil {
			continue
		}
		pages = append(pages, page)
	}
	if len(pages) == 0 {
		pages = make([]models.PageProps, 0)
	}

	ws, err := sh.slackService.CreateWorkspace(
		c.Context(),
		pages,
		resp.IncomingWebhook.ChannelID,
		resp.IncomingWebhook.URL,
		resp.AuthedUser.ID,
		resp.Team.ID,
		resp.AccessToken,
	)
	if err != nil {
		return c.Status(400).SendString(fmt.Sprintf("Error creating workspace: %s", err.Error()))
	}

	sh.logger.Debug("created slack workspace for user", zap.Any("workspace", ws))

	slackDeeplink := fmt.Sprintf("slack://app?team=%s&id=%s&tab=about",
		resp.Team.ID,
		resp.AppID,
	)

	return c.Status(200).JSON(fiber.Map{"deep_link": slackDeeplink})
}

// SlackConfigurationHandler handles the configuration of the Slack app
func (sh *SlackIntegrationHandler) ConfigureCommandHandler(c *fiber.Ctx) error {
	// TODO: remove this
	return c.Status(200).Send(nil)
}

func (sh *SlackIntegrationHandler) WatchCommandHandler(c *fiber.Ctx) error {
	cmd, err := SlashCommandParseFast(c.Request())
	if err != nil {
		return c.Status(400).SendString("Failed to parse command")
	}

	if err := sh.slackService.CreateCompetitorForWorkspace(c.Context(), cmd); err != nil {
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

	if err := sh.slackService.AddUserToSlackWorkspace(c.Context(), cmd); err != nil {
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

	if err := sh.slackService.HandleSlackInteractionPayload(c.Context(), payload); err != nil {
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
