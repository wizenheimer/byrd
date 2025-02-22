package routes

import (
	"github.com/gofiber/fiber/v2"
	handler "github.com/wizenheimer/byrd/src/internal/api/handlers/integration"
	"github.com/wizenheimer/byrd/src/internal/api/middleware"
)

var placeholderHandler = func(c *fiber.Ctx) error { return nil }

func setupIntegrationRoutes(
	app *fiber.App,
	m *middleware.AccessMiddleware,
	sh *handler.SlackIntegrationHandler,
) {
	integration := app.Group("/api/public/v1/integration")

	// List all current integrations
	integration.Get("", placeholderHandler)

	setupSlackIntegrationRoutes(integration, m, sh)
}

func setupSlackIntegrationRoutes(
	router fiber.Router,
	m *middleware.AccessMiddleware,
	sh *handler.SlackIntegrationHandler,
) {
	slack := router.Group("/slack")

	// Handle oauth trigger
	slack.Get("/oauth", sh.SlackOAuthHandler)

	// Handle installation of the slack app to a workspace
	slack.Get("/install", sh.SlackInstallationHandler)

	// Slack command trigger group
	cmdGroup := slack.Group("/cmd", m.RequiresSlackSignature)

	// Handle configure command
	cmdGroup.Post("/configure", sh.ConfigureCommandHandler)

	// Handle watch command
	cmdGroup.Post("/watch", sh.WatchCommandHandler)

	// Handle user command
	cmdGroup.Post("/user", sh.UserCommandHandler)

	// Handle slack app command interactions
	cmdGroup.Post("/interact", sh.SlackInteractionHandler)
}
