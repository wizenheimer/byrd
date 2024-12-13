package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/src/internal/api/handlers"
	"github.com/wizenheimer/iris/src/internal/domain/interfaces"
	"github.com/wizenheimer/iris/src/pkg/logger"
)

type HandlerContainer struct {
	ScreenshotHandler   *handlers.ScreenshotHandler
	DiffHandler         *handlers.DiffHandler
	CompetitorHandler   *handlers.CompetitorHandler
	NotificationHandler *handlers.NotificationHandler
}

func NewHandlerContainer(
	screenshotService interfaces.ScreenshotService,
	diffService interfaces.DiffService,
	competitorService interfaces.CompetitorService,
	notificationService interfaces.NotificationService,
	logger *logger.Logger,
) *HandlerContainer {
	return &HandlerContainer{
		ScreenshotHandler:   handlers.NewScreenshotHandler(screenshotService, logger),
		DiffHandler:         handlers.NewDiffHandler(diffService, logger),
		CompetitorHandler:   handlers.NewCompetitorHandler(competitorService, logger),
		NotificationHandler: handlers.NewNotificationHandler(notificationService, logger),
	}
}

func SetupRoutes(app *fiber.App, handlers *HandlerContainer) {
	api := app.Group("/api/v1")

	// Screenshot routes
	screenshot := api.Group("/screenshot")
	screenshot.Post("/", handlers.ScreenshotHandler.CreateScreenshot)
	screenshot.Get("/image", handlers.ScreenshotHandler.GetScreenshotImage)
	screenshot.Get("/content", handlers.ScreenshotHandler.GetScreenshotContent)

	// Diff routes
	diff := api.Group("/diff")
	diff.Post("/create", handlers.DiffHandler.CreateDiff)
	diff.Get("/report", handlers.DiffHandler.CreateReport)

	// Competitor routes
	competitors := api.Group("/competitors")
	competitors.Post("/", handlers.CompetitorHandler.CreateCompetitor)
	competitors.Get("/", handlers.CompetitorHandler.ListCompetitors)
	competitors.Get("/id/:id", handlers.CompetitorHandler.GetCompetitor)

	// Notification routes
	api.Post("/notify", handlers.NotificationHandler.SendNotification)
}
