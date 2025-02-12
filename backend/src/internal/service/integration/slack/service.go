package slackworkspace

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/slack-go/slack"
	core_models "github.com/wizenheimer/byrd/src/internal/models/core"
	models "github.com/wizenheimer/byrd/src/internal/models/integration/slack"
	repository "github.com/wizenheimer/byrd/src/internal/repository/integration/slack"
	"github.com/wizenheimer/byrd/src/internal/service/workspace"
	"github.com/wizenheimer/byrd/src/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type slackWorkspaceService struct {
	// repo is the repository for Slack workspace data
	repo repository.SlackWorkspaceRepository

	// ws is the workspace service for managing Byrd workspaces
	ws workspace.WorkspaceService

	// logger is the logger for the Slack workspace service
	logger *logger.Logger
}

// NewSlackWorkspaceService creates a new Slack workspace service
func NewSlackWorkspaceService(repo repository.SlackWorkspaceRepository, ws workspace.WorkspaceService, logger *logger.Logger) (SlackWorkspaceService, error) {
	svc := slackWorkspaceService{
		repo:   repo,
		ws:     ws,
		logger: logger,
	}
	return &svc, nil
}

// Creates and associates an existing Byrd workspace with a Slack workspace
func (svc *slackWorkspaceService) CreateSlackWorkspace(ctx context.Context, workspaceID uuid.UUID, teamID, accessToken string) (*models.SlackWorkspace, error) {
	return svc.repo.CreateSlackWorkspace(ctx, workspaceID, teamID, accessToken)
}

// UpdateSlackWorkspace updates the access token for a Slack workspace
func (svc *slackWorkspaceService) UpdateSlackWorkspace(ctx context.Context, cmd slack.SlashCommand) (*models.SlackWorkspace, error) {
	slackWorkspace, err := svc.repo.GetSlackWorkspaceByTeamID(ctx, cmd.TeamID)
	if err != nil {
		return nil, err
	}

	if slackWorkspace.AccessToken == nil {
		return nil, err
	}

	client := slack.New(*slackWorkspace.AccessToken)
	if client == nil {
		return nil, err
	}

	// Get Channel Info
	input := &slack.GetConversationInfoInput{ChannelID: cmd.ChannelID}
	channel, err := client.GetConversationInfo(input)
	if err != nil {
		return nil, err
	}

	var canvasID string
	if channel.Properties != nil && channel.Properties.Canvas.FileId != "" {
		canvasID = channel.Properties.Canvas.FileId
	} else {
		canvasID, err = client.CreateChannelCanvas(cmd.ChannelID, slack.DocumentContent{
			Type:     "markdown",
			Markdown: "Welcome to Byrd! This is your team's shared canvas. Use `/byrd` to interact with Byrd.",
		})
		if err != nil {
			svc.showSupportModal(
				client,
				cmd.TriggerID,
				"Failed to create canvas for channel",
				[]string{
					"Seems like we're having trouble associating the canvas with this channel.",
				})
			return nil, err
		}
	}

	ws, err := svc.repo.UpdateSlackWorkspace(ctx, cmd.TeamID, cmd.ChannelID, canvasID)
	if err != nil {
		svc.showSupportModal(
			client,
			cmd.TriggerID,
			"Failed to associate to workspace repo",
			[]string{
				"Seems like we're having trouble updating the workspace.",
			})
		return nil, err
	}

	// Show Success Modal
	svc.showSuccessModal(
		client,
		cmd.TriggerID,
		cmd.ChannelID,
		"",
		"Your slack channel is now in sync with Byrd.",
		[]string{
			"`/watch` adds a new page to your watchlist.",
			"`/invite` lets you bring your team along.",
			"And that's it! you're all set to go.",
		},
	)

	return ws, nil
}

// Handles the bookkeeping for a Slack integration that has been removed
func (svc *slackWorkspaceService) DeleteSlackWorkspace(ctx context.Context, workspaceID uuid.UUID) error {
	ws, err := svc.repo.GetSlackWorkspaceByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return err
	}
	return svc.repo.DeleteSlackWorkspace(ctx, ws.TeamID)
}

// GetSlackWorkspaceByTeamID retrieves a Slack workspace by its team ID
func (svc *slackWorkspaceService) GetSlackWorkspaceByTeamID(ctx context.Context, teamID string) (*models.SlackWorkspace, error) {
	return svc.repo.GetSlackWorkspaceByTeamID(ctx, teamID)
}

// GetSlackWorkspaceByWorkspaceID retrieves a Slack workspace by its Byrd workspace ID
func (svc *slackWorkspaceService) GetSlackWorkspaceByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (*models.SlackWorkspace, error) {
	return svc.repo.GetSlackWorkspaceByWorkspaceID(ctx, workspaceID)
}

// BatchGetWorkspaceByWorkspaceID retrieves a list of Slack workspaces by their Byrd workspace IDs
func (svc *slackWorkspaceService) BatchGetWorkspaceByWorkspaceIDs(ctx context.Context, workspaceIDs []uuid.UUID) ([]*models.SlackWorkspace, error) {
	return svc.repo.BatchGetSlackWorkspacesByWorkspaceIDs(ctx, workspaceIDs)
}

// IntegrationExistsForWorkspace checks if a Slack integration exists for a Byrd workspace
func (svc *slackWorkspaceService) IntegrationExistsForWorkspace(ctx context.Context, workspaceID uuid.UUID) (bool, error) {
	ws, err := svc.repo.GetSlackWorkspaceByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return false, err
	}

	return ws != nil, nil
}

// ------ USER MANAGEMENT ------ //

// AddUserToSlackWorkspace adds a user to a Slack workspace
func (svc *slackWorkspaceService) AddUserToSlackWorkspace(ctx context.Context, cmd slack.SlashCommand) error {
	creatorID := cmd.UserID

	modalRequest := slack.ModalViewRequest{
		Type:       "modal",
		Title:      slack.NewTextBlockObject(slack.PlainTextType, "Invite Users", false, false),
		Submit:     slack.NewTextBlockObject(slack.PlainTextType, "Send Invites", false, false),
		Close:      slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false),
		CallbackID: "invite_users",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewInputBlock(
					"user_selection",
					slack.NewTextBlockObject(slack.PlainTextType, "Users", false, false),
					nil,
					slack.NewOptionsSelectBlockElement(
						slack.OptTypeUser,
						slack.NewTextBlockObject(slack.PlainTextType, "Select users", false, false),
						"user_select",
					),
				),
			},
		},
		PrivateMetadata: creatorID,
	}

	ws, err := svc.repo.GetSlackWorkspaceByTeamID(ctx, cmd.TeamID)
	if err != nil {
		return err
	}

	if ws.AccessToken == nil {
		return errors.New("no access token found for Slack workspace")
	}

	client := slack.New(*ws.AccessToken)

	canAddUsers, workspacePlan, err := svc.ws.CanAddUsers(ctx, ws.WorkspaceID, 1)
	if err != nil {
		svc.showSupportModal(
			client,
			cmd.TriggerID,
			"Couldn't add user",
			[]string{
				"Seems like we're having trouble adding a user.",
			},
		)
		return nil
	}

	if !canAddUsers {
		if err := svc.showUsageLimitModal(
			client,
			cmd.TriggerID,
			workspacePlan,
			core_models.WorkspaceResourceUsers,
		); err != nil {
			svc.showSupportModal(
				client,
				cmd.TriggerID,
				"Couldn't add user",
				[]string{
					"Seems like we're having trouble adding a user.",
				},
			)
		}
		return nil
	}

	_, err = client.OpenView(cmd.TriggerID, modalRequest)
	if err != nil {
		return err
	}

	return nil
}

// UserExistsInSlackWorkspace checks if a user exists in a Slack workspace
func (svc *slackWorkspaceService) UserExistsInSlackWorkspace(ctx context.Context, teamID string, userEmail string) (bool, error) {
	return false, nil
}

// ------ COMPETITOR MANAGEMENT ------ //

// CreateCompetitor creates a competitor in a Slack workspace
func (svc *slackWorkspaceService) CreateCompetitorForWorkspace(ctx context.Context, cmd slack.SlashCommand) error {
	args := strings.Fields(cmd.Text) // Extract URLs from command
	if len(args) == 0 {
		return errors.New("no URLs provided")
	}

	ws, err := svc.repo.GetSlackWorkspaceByTeamID(ctx, cmd.TeamID)
	if err != nil {
		return err
	}

	if ws.AccessToken == nil {
		return errors.New("no access token found for Slack workspace")
	}

	client := slack.New(*ws.AccessToken)

	var urlBlocks []slack.Block
	var urls []string
	for _, u := range args {
		url, err := url.Parse(u) // Ensure URL is valid
		if err != nil {
			continue
		}
		urls = append(urls, url.String())
	}

	// Check if the user can create a page
	canCreatePage, workspacePlan, err := svc.ws.CanCreatePage(ctx, ws.WorkspaceID, len(urls))
	if err != nil {
		svc.showSupportModal(
			client,
			cmd.TriggerID,
			"Failed to create page",
			[]string{
				"Seems like we're having trouble creating a page.",
			},
		)
		return nil
	}
	if !canCreatePage {
		if err := svc.showUsageLimitModal(
			client,
			cmd.TriggerID,
			workspacePlan,
			core_models.WorkspaceResourcePages,
		); err != nil {
			svc.showSupportModal(
				client,
				cmd.TriggerID,
				"Failed to show usage limit modal",
				[]string{
					"Seems like we're having trouble creating a page.",
				},
			)
		}
		return nil
	}

	// Check if the user can create a competitor
	canCreateCompetitor, workspacePlan, err := svc.ws.CanCreateCompetitor(ctx, ws.WorkspaceID, 1, len(urls))
	if err != nil {
		svc.showSupportModal(
			client,
			cmd.TriggerID,
			"Failed to create competitor",
			[]string{
				"Seems like we're having trouble creating a competitor.",
			},
		)
		return nil
	}
	if !canCreateCompetitor {
		if err := svc.showUsageLimitModal(
			client,
			cmd.TriggerID,
			workspacePlan,
			core_models.WorkspaceResourceCompetitors,
		); err != nil {
			svc.showSupportModal(
				client,
				cmd.TriggerID,
				"Failed to show usage limit modal",
				[]string{
					"Seems like we're having trouble creating a page.",
				},
			)
		}
		return nil
	}

	// --- Dropdown (Single Select) ---
	competitorSelect := slack.NewOptionsSelectBlockElement(
		slack.OptTypeStatic,
		slack.NewTextBlockObject(slack.PlainTextType, "Select a competitor", false, false),
		"select_competitor",
		svc.getCompetitorOptions(ctx, ws.WorkspaceID)...,
	)

	competitorBlock := slack.NewInputBlock(
		"competitor_selection",
		slack.NewTextBlockObject(slack.PlainTextType, "Assign to Competitor", false, false),
		nil, // No hint
		competitorSelect,
	)

	// --- Multi-Select (DiffProfile) ---
	caser := cases.Title(language.English)
	diffProfileOptions := core_models.GetDefaultDiffProfile()
	var multiSelectOptions []*slack.OptionBlockObject
	for _, profile := range diffProfileOptions {
		multiSelectOptions = append(multiSelectOptions, slack.NewOptionBlockObject(
			profile,
			slack.NewTextBlockObject(slack.PlainTextType, caser.String(profile), false, false),
			nil,
		))
	}

	diffProfileMultiSelect := slack.NewOptionsMultiSelectBlockElement(
		slack.MultiOptTypeStatic,
		slack.NewTextBlockObject(slack.PlainTextType, "Product, Pricing, Partnerships etc.", false, false),
		"select_diff_profiles",
		multiSelectOptions...,
	)

	diffProfileBlock := slack.NewInputBlock(
		"diff_profile_selection",
		slack.NewTextBlockObject(slack.PlainTextType, "Select Competitor Profiles", false, false),
		nil, // No hint
		diffProfileMultiSelect,
	)

	competitorData := CompetitorData{
		ChannelID: cmd.ChannelID,
		URLs:      urls,
	}

	jsonBytes, err := json.Marshal(competitorData)
	if err != nil {
		return err
	}
	base64String := base64.StdEncoding.EncodeToString(jsonBytes)

	modal := slack.ModalViewRequest{
		Type:            slack.VTModal,
		Title:           slack.NewTextBlockObject(slack.PlainTextType, "Assign Competitor", false, false),
		Submit:          slack.NewTextBlockObject(slack.PlainTextType, "Save", false, false),
		Close:           slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false),
		Blocks:          slack.Blocks{BlockSet: append(urlBlocks, competitorBlock, diffProfileBlock)}, // Ensure correct structure
		CallbackID:      "save_competitor",
		PrivateMetadata: base64String,
	}

	_, err = client.OpenView(cmd.TriggerID, modal)
	if err != nil {
		return err
	}

	return nil
}

type CompetitorData struct {
	ChannelID string   `json:"channel_id"`
	URLs      []string `json:"urls"`
}

func (svc *slackWorkspaceService) getCompetitorOptions(ctx context.Context, workspaceID uuid.UUID) []*slack.OptionBlockObject {
	competitors, _, err := svc.ws.ListCompetitorsForWorkspace(ctx, workspaceID, nil, nil)
	if err != nil {
		return nil
	}

	options := make([]*slack.OptionBlockObject, 0)

	newCompetitorOption := slack.NewOptionBlockObject(
		uuid.Nil.String(),
		slack.NewTextBlockObject(
			slack.PlainTextType,
			"Create New Competitor",
			false,
			false,
		), nil)
	options = append(options, newCompetitorOption)

	for _, competitor := range competitors {
		competitorOption := slack.NewOptionBlockObject(
			competitor.ID.String(),
			slack.NewTextBlockObject(
				slack.PlainTextType,
				competitor.Name,
				false,
				false,
			), nil)
		options = append(options, competitorOption)
	}

	return options
}

// AddPageToCompetitor adds a page to a competitor in a Slack workspace
func (svc *slackWorkspaceService) AddPageToCompetitor(ctx context.Context, teamID string, competitorID uuid.UUID, pageURLs []string) error {
	return nil
}

func (svc *slackWorkspaceService) HandleSlackInteractionPayload(ctx context.Context, payload slack.InteractionCallback) error {

	switch payload.Type {
	case slack.InteractionTypeViewSubmission:
		payloadCallback := payload.View.CallbackID
		switch payloadCallback {
		case "save_competitor":
			return svc.handlePageSubmission(ctx, payload)
		case "support_submission":
			return svc.handleSupportSubmission(ctx, payload)
		case "invite_users":
			return svc.handleInviteSubmission(ctx, payload)
		}
	case slack.InteractionTypeBlockActions:
		return svc.handleInviteResponse(ctx, payload)
	}

	return errors.New("unsupported interaction type")
}

// TODO: remove this function
func (svc *slackWorkspaceService) getUserEmail(client *slack.Client, userID string) (string, error) {
	user, err := client.GetUserInfo(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %w", err)
	}

	return user.Profile.Email, nil
}

func (svc *slackWorkspaceService) handleInviteSubmission(ctx context.Context, payload slack.InteractionCallback) error {
	selectedUser := payload.View.State.Values["user_selection"]["user_select"].SelectedUser
	if selectedUser == "" {
		return errors.New("no user selected")
	}

	inviteCreator := payload.View.PrivateMetadata // Retrieve stored creator ID

	svc.logger.Info("User selected for invite", zap.String("user_id", selectedUser), zap.String("creator_id", inviteCreator))

	// Retrieve Slack workspace token
	teamID := payload.User.TeamID
	ws, err := svc.repo.GetSlackWorkspaceByTeamID(ctx, teamID)
	if err != nil {
		return err
	}
	if ws.AccessToken == nil {
		return errors.New("no access token found for Slack workspace")
	}

	client := slack.New(*ws.AccessToken)

	// Open a DM with the selected user
	channel, _, _, err := client.OpenConversation(&slack.OpenConversationParameters{
		Users: []string{selectedUser},
	})
	if err != nil {
		svc.logger.Error("failed to open DM with user", zap.Error(err))
		return err
	}

	// Construct the invite message with action buttons
	msgBlocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject(slack.MarkdownType, fmt.Sprintf(
				"*%s* invited you to join the workspace! 🎉\nWould you like to accept or decline?",
				fmt.Sprintf("<@%s>", inviteCreator), // Mention creator
			), false, false),
			nil, nil,
		),
		slack.NewActionBlock("invite_response",
			slack.NewButtonBlockElement("accept_invite", inviteCreator, slack.NewTextBlockObject(slack.PlainTextType, "✅ Accept", false, false)),
			slack.NewButtonBlockElement("decline_invite", inviteCreator, slack.NewTextBlockObject(slack.PlainTextType, "❌ Decline", false, false)),
		),
	}

	// Send the message as a DM
	_, _, err = client.PostMessage(channel.ID, slack.MsgOptionBlocks(msgBlocks...))
	if err != nil {
		svc.logger.Error("failed to send invite message", zap.Error(err))
		return err
	}

	return nil
}

func (svc *slackWorkspaceService) handleInviteResponse(ctx context.Context, payload slack.InteractionCallback) error {
	actionID := payload.ActionCallback.BlockActions[0].ActionID
	userID := payload.User.ID
	// Retrieve workspace token
	teamID := payload.Team.ID
	ws, err := svc.repo.GetSlackWorkspaceByTeamID(ctx, teamID)
	if err != nil {
		return err
	}
	if ws.AccessToken == nil {
		return errors.New("no access token found for Slack workspace")
	}

	client := slack.New(*ws.AccessToken)

	var responseMessage string

	if actionID == "accept_invite" {
		responseMessage = "And just like that, you're in! 🎉"

		userEmail, err := svc.getUserEmail(client, payload.User.ID)
		if err != nil {
			responseMessage = "Seems like we're having trouble getting your email address."
			svc.logger.Error("failed to get user email", zap.Error(err))
		}

		inviteCreator := payload.ActionCallback.BlockActions[0].Value // Get creator ID from button value

		inviteCreatorEmail, err := svc.getUserEmail(client, inviteCreator)
		if err != nil {
			responseMessage = "Seems like we're having trouble getting the creator's email address."
			svc.logger.Error("failed to get invite creator email", zap.Error(err))
		}

		if inviteCreatorEmail != "" && userEmail != "" {
			workspaceUsers, err := svc.ws.AddSlackUserToWorkspace(
				ctx,
				inviteCreatorEmail,
				ws.WorkspaceID,
				[]string{
					userEmail,
				},
			)
			if err != nil {
				responseMessage = "Seems like we're having trouble adding you to the workspace."
				svc.logger.Error("failed to add user to workspace", zap.Error(err))
			}
			if len(workspaceUsers) == 0 {
				responseMessage = "Seems you tried to join a workspace you're already a part of. That's cool too!"
			}
		}

	} else if actionID == "decline_invite" {
		responseMessage = "No worries! We are ready when you are. 🙌"
	} else {
		return errors.New("wait, how did you get here?")
	}

	// Send a confirmation message to the user
	_, _, err = client.PostMessage(userID, slack.MsgOptionText(responseMessage, false))
	if err != nil {
		svc.logger.Error("failed to send response message", zap.Error(err))
		return err
	}

	return nil
}
