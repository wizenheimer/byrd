package routes

import (
	"github.com/gofiber/fiber/v2"
	intg_handler "github.com/wizenheimer/byrd/src/internal/api/handlers/integration"
)

var placeholderHandler = func(c *fiber.Ctx) error { return nil }

func setupIntegrationRoutes(
	app *fiber.App,
	sh *intg_handler.SlackIntegrationHandler,
) {
	integration := app.Group("/api/public/v1/integration")

	// List all current integrations
	integration.Get("", placeholderHandler)

	setupSlackIntegrationRoutes(integration, sh)
}

func setupSlackIntegrationRoutes(
	router fiber.Router,
	sh *intg_handler.SlackIntegrationHandler,
) {
	slack := router.Group("/slack")

	// Handle oauth trigger
	slack.Get("/oauth", sh.SlackOAuthHandler)

	// Handle installation of the slack app to a workspace
	slack.Get("/install", sh.SlackInstallationHandler)

	// Slack command trigger group
	cmdGroup := slack.Group("/cmd")

	// Handle configure command
	cmdGroup.Post("/configure", sh.ConfigureCommandHandler)

	// Handle watch command
	cmdGroup.Post("/watch", sh.WatchCommandHandler)

	// Handle user command
	cmdGroup.Post("/user", sh.UserCommandHandler)

	// Handle slack app command interactions
	cmdGroup.Post("/interact", sh.SlackInteractionHandler)

	// Slack event group
	eventGroup := slack.Group("/event")

	// Handle slack app events
	eventGroup.Post("", placeholderHandler)
}
