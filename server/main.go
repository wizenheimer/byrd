package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/wizenheimer/iris/internal/api/middleware"
	"github.com/wizenheimer/iris/internal/api/routes"
	"github.com/wizenheimer/iris/internal/config"
	"github.com/wizenheimer/iris/internal/domain/interfaces"
	"github.com/wizenheimer/iris/internal/domain/models"
	"github.com/wizenheimer/iris/internal/repository/db"
	"github.com/wizenheimer/iris/internal/repository/storage"
	"github.com/wizenheimer/iris/internal/service/ai"
	"github.com/wizenheimer/iris/internal/service/competitor"
	"github.com/wizenheimer/iris/internal/service/diff"
	"github.com/wizenheimer/iris/internal/service/notification"
	"github.com/wizenheimer/iris/internal/service/screenshot"
	"github.com/wizenheimer/iris/pkg/utils/sqlDb"
)

func initializer(cfg *config.Config, sqlDb *sql.DB) (*routes.HandlerContainer, error) {
	// Initialize HTTP client
	httpClient := &http.Client{}

	// Initialize services
	screenshotService, err := setupScreenshotService(cfg, httpClient)
	if err != nil {
		return nil, err
	}

	notificationService, err := setupNotificationService(cfg, httpClient)
	if err != nil {
		return nil, err
	}

	diffService, err := setupDiffService(cfg, sqlDb, httpClient)
	if err != nil {
		return nil, err
	}

	competitorService, err := setupCompetitorService(sqlDb)
	if err != nil {
		return nil, err
	}

	// Initialize handlers
	handlers := routes.NewHandlerContainer(
		screenshotService,
		diffService,
		competitorService,
		notificationService,
	)

	return handlers, nil
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	// Initialize database
	sqlDb, err := sqlDb.Init(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
		return
	}
	defer sqlDb.Close()

	// Initialize handlers
	handlers, err := initializer(cfg, sqlDb)
	if err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
		return
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: middleware.CustomErrorHandler,
	})

	// Setup middleware
	middleware.SetupMiddleware(app)

	// Setup routes
	routes.SetupRoutes(app, handlers)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(cfg.Server.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // This blocks the main thread until an interrupt is received
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to gracefully shutdown the server: %v", err)
	}
}

func setupScreenshotService(cfg *config.Config, httpClient interfaces.HTTPClient) (interfaces.ScreenshotService, error) {
	storageRepo, err := storage.NewR2Storage(
		cfg.Storage.AccessKey,
		cfg.Storage.SecretKey,
		cfg.Storage.Bucket,
		cfg.Storage.Region,
	)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	return screenshot.NewScreenshotService(
		storageRepo,
		httpClient,
		&models.ScreenshotServiceConfig{},
	)
}

func setupNotificationService(cfg *config.Config, httpClient *http.Client) (interfaces.NotificationService, error) {
	emailClient, err := notification.NewResendEmailClient(cfg, httpClient)
	if err != nil {
		log.Fatalf("Failed to initialize email client: %v", err)
	}

	templateManager, err := notification.NewTemplateManager()
	if err != nil {
		log.Fatalf("Failed to initialize template manager: %v", err)
	}

	return notification.NewNotificationService(emailClient, templateManager)
}

func setupDiffService(cfg *config.Config, sqlDb *sql.DB, httpClient *http.Client) (interfaces.DiffService, error) {
	screenshotService, err := setupScreenshotService(cfg, httpClient)
	if err != nil {
		log.Fatalf("Failed to initialize screenshot service: %v", err)
	}

	aiService, err := ai.NewOpenAIService(cfg.Services.OpenAIKey, httpClient)
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}

	diffRepo, err := db.NewDiffRepository(sqlDb)
	if err != nil {
		log.Fatalf("Failed to initialize diff repository: %v", err)
	}

	return diff.NewDiffService(diffRepo, aiService, screenshotService)
}

func setupCompetitorService(sqlDb *sql.DB) (interfaces.CompetitorService, error) {
	competitorRepo, err := db.NewCompetitorRepository(sqlDb)
	if err != nil {
		log.Fatalf("Failed to initialize competitor repository: %v", err)
	}

	return competitor.NewCompetitorService(competitorRepo)
}
