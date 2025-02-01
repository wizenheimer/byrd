// ./src/server/startup/client.go
package startup

import (
	"github.com/wizenheimer/byrd/src/internal/client"
	"github.com/wizenheimer/byrd/src/internal/config"
	"github.com/wizenheimer/byrd/src/pkg/logger"
)

func SetupScreenshotClient(cfg *config.Config, logger *logger.Logger) (*client.HTTPClient, error) {
	screenshotClientOpts := []client.ClientOption{
		client.WithRateLimit(cfg.Services.ScreenshotServiceQPS, 1),
		client.WithRetry(3, []int{408, 429, 500, 502, 503, 504}),
	}

	client, err := client.NewClient(logger, screenshotClientOpts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}
