package template

import (
	"fmt"
	"time"
)

type TemplateName string

const (
	// -- workspace user lifecycle templates --
	WorkspaceInvitePendingTemplate  TemplateName = "pending_workspace_invite"
	WorkspaceInviteAcceptedTemplate TemplateName = "accept_workspace_invite"

	// -- weekly roundup templates --
	WeeklyRoundupTemplate TemplateName = "weekly_roundup"
)

var templates = map[TemplateName]Template{
	WorkspaceInvitePendingTemplate:  workspaceInvitePending,
	WorkspaceInviteAcceptedTemplate: workspaceInviteAccepted,
	WeeklyRoundupTemplate:           weeklyRoundup,
}

// registerDefaultTemplates pre-registers all the default email templates
func registerDefaultTemplates(lib TemplateLibrary) error {
	// Register each template
	for name, tmpl := range templates {
		if err := lib.RegisterTemplate(name, tmpl); err != nil {
			return fmt.Errorf("failed to register template %s: %w", name, err)
		}
	}

	return nil
}

var (
	// WorkspaceInvitePendingTemplate is the template for a pending workspace invite
	workspaceInvitePending = &CommonTemplate{
		PreviewText: "Welcome to the Team!",
		Title:       "You just got added to a workspace",
		Body: []string{
			"Hi there,",
			"Your team's been busy in here! Come Join Your Crew!",
		},
		CTA: &CallToAction{
			ButtonText: "Join Workspace",
			ButtonURL:  "https://byrdhq.com/dashboard/invites",
		},
		ClosingText: "Secure by design - no password needed.",
		Footer: Footer{
			ContactMessage: "Need help? We've got your back:",
			ContactEmail:   "hey@byrd.com",
		},
		GeneratedAt: time.Now(),
	}

	// WorkspaceInviteAcceptedTemplate is the template for an accepted workspace invite
	workspaceInviteAccepted = &CommonTemplate{
		PreviewText: "You're In! Own Their Next Move Before They Make It",
		Title:       "Turn Their Next Move Into Your Next Win",
		Body: []string{
			"Hi there,",
			"Wish you could see every move your competitors make? Well, now you can. With Byrd, every change, no matter how small, gets flagged. From pricing shifts to product updates, you'll be the first to know.",
		},
		BulletTitle: "What's Included:",
		Bullets: []string{
			"Page Monitoring - Stay ahead of every product move",
			"Inbox Monitoring - Monitor the direct line to their customers",
			"Social Monitoring - Keep a pulse on their community",
			"Review Monitoring - They churn. You learn.",
		},
		CTA: &CallToAction{
			ButtonText: "Get Started Now",
			ButtonURL:  "https://byrdhq.com/dashboard",
		},
		Footer: Footer{
			ContactMessage: "Need help? We've got your back:",
			ContactEmail:   "hey@byrd.com",
		},
		GeneratedAt: time.Now(),
	}

	// WeeklyRoundupTemplate is the template for a weekly roundup
	weeklyRoundup = &SectionedTemplate{}
)
