package utils

import "github.com/robfig/cron/v3"

func NewScheduleParser() cron.Parser {
	parser := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

	return parser
}
