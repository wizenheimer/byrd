package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/byrd/src/internal/api/handlers"
)

var placeholderHandler = func(c *fiber.Ctx) error { return nil }

func setupIntegrationRoutes(
	app *fiber.App,
	sh *handlers.SlackIntegrationHandler,
) {
	integration := app.Group("/api/public/v1/integration")

	// List all current integrations
	integration.Get("", placeholderHandler)

	setupSlackIntegrationRoutes(integration, sh)
}

func setupSlackIntegrationRoutes(
	router fiber.Router,
	sh *handlers.SlackIntegrationHandler,
) {
	slack := router.Group("/slack")

	// Handle installation of the slack app to a workspace
	slack.Post("", sh.SlackInstallationHandler)

	// Handle uninstallation of the slack app from a workspace
	slack.Delete("", sh.SlackUninstallationHandler)

	// Handle slack app status check
	slack.Get("", sh.StatusSlackHandler)

	// Handle oauth redirect
	slack.Get("/oauth", sh.SlackOAuthHandler)

	// Slack command trigger group
	cmdGroup := slack.Group("/cmd")

	// Handle slack app command triggers
	cmdGroup.Post("", placeholderHandler)

	// Handle slack app command interactions
	cmdGroup.Post("/interact", placeholderHandler)

	// Slack event group
	eventGroup := slack.Group("/event")

	// Handle slack app events
	eventGroup.Post("", placeholderHandler)
}
