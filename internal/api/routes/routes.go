package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/internal/api/handlers"
	"github.com/wizenheimer/iris/internal/domain/interfaces"
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
) *HandlerContainer {
	return &HandlerContainer{
		ScreenshotHandler:   handlers.NewScreenshotHandler(screenshotService),
		DiffHandler:         handlers.NewDiffHandler(diffService),
		CompetitorHandler:   handlers.NewCompetitorHandler(competitorService),
		NotificationHandler: handlers.NewNotificationHandler(notificationService),
	}
}

func SetupRoutes(app *fiber.App, handlers *HandlerContainer) {
	api := app.Group("/api/v1")

	// Screenshot routes
	screenshot := api.Group("/screenshot")
	screenshot.Post("/", handlers.ScreenshotHandler.CreateScreenshot)
	screenshot.Get("/:hash/:weekNumber/:weekDay", handlers.ScreenshotHandler.GetScreenshot)
	screenshot.Get("/content/:hash/:weekNumber/:weekDay", handlers.ScreenshotHandler.GetScreenshotContent)

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
