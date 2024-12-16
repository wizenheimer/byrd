package interfaces

import "context"

// SubscriptionRepository defines subscription capabilities
type SubscriptionRepository interface {
	// Subscribe subscribes an email to a competitor
	Subscribe(ctx context.Context, competitorID int, email string) error
	// Unsubscribe unsubscribes an email from a competitor
	Unsubscribe(ctx context.Context, competitorID int, email string) error
	// GetSubscribersByCompetitor gets all subscribers for a competitor
	GetSubscribersByCompetitor(ctx context.Context, competitorID int) ([]string, error)
}
