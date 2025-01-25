// ./src/server/startup/client.go
package startup

import (
	"github.com/wizenheimer/byrd/src/internal/client"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

func SetupScreenshotClient(cfg *config.Config, logger *logger.Logger) (client.HTTPClient, error) {
	screenshotClientOpts := []client.ClientOption{
		client.WithLogger(logger),
		client.WithAuth(client.BearerAuth{
			Token: cfg.Services.ScreenshotServiceAPIKey,
		}),
	}

	screenshotHttpClient, err := client.NewClient(screenshotClientOpts...)
	if err != nil {
		return nil, err
	}

	return client.NewRateLimitedClient(screenshotHttpClient, cfg.Services.ScreenshotServiceQPS), nil
}
