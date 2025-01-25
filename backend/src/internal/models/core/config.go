package models

// SlackConfig holds Slack-specific configuration
type SlackConfig struct {
	Token        string
	ChannelID    string
	DefaultUser  string
	RetryCount   int
	ColorMapping map[Severity]string
}

// DefaultSlackConfig returns a default configuration
func DefaultSlackConfig() SlackConfig {
	return SlackConfig{
		RetryCount:  3,
		DefaultUser: "Byrd Dev",
		ColorMapping: map[Severity]string{
			SeverityInfo:     "#36a64f",
			SeverityWarning:  "#ffaa00",
			SeverityError:    "#ff0000",
			SeverityCritical: "#7b0000",
		},
	}
}
