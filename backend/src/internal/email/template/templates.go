package template

import "time"

// Event driven email templates
// These emails are triggered by specific events in the user lifecycle
// And expect realtime delivery

// welcomeTemplate is a common template for welcome emails
// Event: Triggered when a user signs up for a trial and is granted access
var welcomeTemplate = CommonTemplate{
	PreviewText: "Your Competition's Worst Nightmare Just Got Real",
	Title:       "Your Competition's Worst Nightmare Just Got Real",
	Subtitle:    "The wait is over - We are officially yours!",
	Body: []string{
		"Hi there,",
		"Wish you could see every move your competitors make? Well, now you can. With Byrd, every change, no matter how small, gets flagged. From pricing shifts to product updates, you'll be the first to know.",
	},
	BulletTitle: "WHAT'S INCLUDED",
	Bullets: []string{
		"Page Monitoring - Stay ahead of every product move",
		"Inbox Monitoring - Monitor the direct line to their customers",
		"Social Monitoring - Keep a pulse on their community",
		"Review Monitoring - They churn. You learn.",
	},
	CTA: &CallToAction{
		ButtonText: "Get Started Now",
		ButtonURL:  "https://byrdhq.com/dashboard",
		FooterText: "They're hoping this invite expires. Disappoint them.",
	},
	Footer: Footer{
		ContactMessage: "Need help? We've got your back:",
		ContactEmail:   "hey@byrd.com",
	},
	GeneratedAt: time.Now(),
}

// trialWaitlistTemplate is a common template for trial waitlist emails
// Event: Triggered when a user signs up for a trial and is placed on a waitlist
var waitlistTemplate = CommonTemplate{
	PreviewText: "Good things come to those who... actually, let's speed this up",
	Title:       "You're Almost There",
	Body: []string{
		"Hi there,",
		"We're thrilled to have you on board! We're gradually rolling out access for new teams, and will reach out to you with your onboarding details as soon as your spot opens up.",
		"While we know waiting isn't ideal, we've prepared some excellent resources to help you make the most of this time. From swipe files to sales battlecards, you'll have access to the exact tools our top users rely on - and they're yours to keep!",
	},
	ClosingText: "Can't Wait? Don't Wait\nIncase waiting doesn't work for you (we understand!), reach out to us. We're founders too and we occasionally fast-track access for teams who are ready to dive straight in.",
	Footer: Footer{
		ContactMessage: "Need help? We've got your back:",
		ContactEmail:   "hey@byrd.com",
	},
	GeneratedAt: time.Now(),
}

// pendingWorkspaceInviteTemplate is a common template for pending workspace invite emails
// Event: Triggered when a user is invited to a workspace
var pendingWorkspaceInviteTemplate = CommonTemplate{
	PreviewText: "Welcome to the Team!",
	Title:       "Sarah Parker added you to their workspace",
	Body: []string{
		"Hi there,",
		"Sarah Parker brought you into your team's workspace. Come Join Your Crew!",
	},
	CTA: &CallToAction{
		ButtonText: "Join Workspace",
		ButtonURL:  "https://byrdhq.com/dashboard",
	},
	ClosingText: "Secure by design - no password needed.",
	GeneratedAt: time.Now(),
}

// acceptWorkspaceInviteTemplateFTU is a common template for accept workspace invite emails
// This is a First Time User (FTU) template
// Event: Triggered when a user accepts an invite to a workspace
var acceptWorkspaceInviteTemplateFTU = CommonTemplate{
	PreviewText: "You're In! Own Their Next Move Before They Make It",
	Title:       "Turn Their Next Move Into Your Next Win",
	Body: []string{
		"Hi there,",
		"Wish you could see every move your competitors make? Well, now you can. With Byrd, every change, no matter how small, gets flagged. From pricing shifts to product updates, you'll be the first to know.",
	},
	BulletTitle: "WHAT'S INCLUDED",
	Bullets: []string{
		"Page Monitoring - Stay ahead of every product move",
		"Inbox Monitoring - Monitor the direct line to their customers",
		"Social Monitoring - Keep a pulse on their community",
		"Review Monitoring - They churn. You learn.",
	},
	CTA: &CallToAction{
		ButtonText: "Get Started Now",
		ButtonURL:  "https://byrdhq.com/dashboard",
		FooterText: "They're hoping this invite expires. Disappoint them.",
	},
	Footer: Footer{
		ContactMessage: "Need help? We've got your back:",
		ContactEmail:   "hey@byrd.com",
	},
	GeneratedAt: time.Now(),
}

// acceptWorkspaceInviteTemplate is a common template for accept workspace invite emails
// This is a Returning User (RU) template
// Event: Triggered when a user accepts an invite to a workspace
var acceptWorkspaceInviteTemplateRU = CommonTemplate{
	PreviewText: "You're In! Own Their Next Move Before They Make It",
	Title:       "Because Knowing First Means Moving First",
	Body: []string{
		"Hi there,",
		"Great to have you here! Your workspace is ready to help you track updates that matters. No more missing important updates or scrambling for information.",
	},
	CTA: &CallToAction{
		ButtonText: "Get started",
		ButtonURL:  "https://byrdhq.com/dashboard",
	},
	ClosingText: "Quick win: Link your Slack now and watch the magic happen!",
	GeneratedAt: time.Now(),
}

// declineWorkspaceInviteTemplate is a common template for decline workspace invite emails
// Event: Triggered when a user declines an invite to a workspace
var declineWorkspaceInviteTemplate = CommonTemplate{
	PreviewText: "Quick note about your invite",
	Title:       "We're keeping your spot warm, just in case",
	Body: []string{
		"Hi there,",
		"Thanks for letting us know about the workspace invite. We'll keep your spot warm, just in case you change your mind.",
		"Want to join the workspace later? We'll pass your request to the workspace owner.",
	},
	CTA: &CallToAction{
		ButtonText: "Request Invite",
		ButtonURL:  "https://byrdhq.com/request-invite",
	},
	ClosingText: "No hard feelings. We're here when you're ready.",
	GeneratedAt: time.Now(),
}

// requestWorkspaceInviteTemplate is a common template for request workspace invite emails
// Event: Triggered when a user requests an invite to a workspace
var requestWorkspaceInviteTemplate = CommonTemplate{
	PreviewText: "Someone wants to join your team, Sarah Parker",
	Title:       "You've got a request to join your workspace",
	Body: []string{
		"Hi there,",
		"Sarah Parker wants to join your workspace. Let's get them in!",
	},
	CTA: &CallToAction{
		ButtonText: "Approve Request",
		ButtonURL:  "https://byrdhq.com/workspace-invites",
	},
	ClosingText: "Winning is a team sport. Let's get your heavy hitters!",
	GeneratedAt: time.Now(),
}

// trialSucceededTemplate is a common template for trial succeeded emails
// Event: Triggered when billing is successful for a trial user
var trialSucceededTemplate = CommonTemplate{
	PreviewText: "You made our nights and weekends worth it (and yes, your access is confirmed)",
	Title:       "More Than Just a Renewal",
	Subtitle:    "Team just did a happy dance (your renewal triggered it)",
	Body: []string{
		"Hi there,",
		"We know this is supposed to be a standard payment confirmation email, but instead we wanted to drop by and say a heartfelt thanks - not just for the wire, but for all the feedback that's helped make Byrd better.",
		"Teams like yours are why we get excited about competitive intelligence. To make this better, we've snuck a few extra user seats to your account (on the house). Here's to more winning moves, together with Byrd.",
	},
	Footer:      Footer{},
	GeneratedAt: time.Now(),
}

// deletedWorkspaceTemplate is a common template for deleted workspace emails
// Event: Triggered when a workspace is deleted
var deletedWorkspaceTemplate = CommonTemplate{
	PreviewText: "Quick update about your workspace",
	Title:       "Your workspace has been deleted",
	Body: []string{
		"Hi there,",
		"Just confirming that your workspace has been deleted. Everything's been cleared out safely.",
		"Changed your mind? We can help you start fresh anytime.",
	},
	CTA: &CallToAction{
		ButtonText: "Create Workspace",
		ButtonURL:  "https://byrdhq.com/create-workspace",
	},
	ClosingText: "We're here to help if you need anything.",
	GeneratedAt: time.Now(),
}

// renewalFailedTemplate is a common template for renewal failed emails
// Event: Triggered when a renewal payment fails
var renewalFailedTemplate = CommonTemplate{
	PreviewText: "Quick note on your renewal",
	Title:       "Your renewal payment didn't go through (but don't worry, you're still good)",
	Body: []string{
		"Hi there,",
		"We noticed your renewal payment didn't go through. To make sure you don't lose access to your workspace, we've extended your access for another 7 days. We'll try to charge your card again in a week.",
		"If you'd like to take a look at what might be happening on your end, we're here to help",
	},
	CTA: &CallToAction{
		ButtonText: "Update Payment",
		ButtonURL:  "https://byrdhq.com/dashboard",
		FooterText: "Maybe it's just a small hiccup. Let's get this sorted.",
	},
	Footer: Footer{
		ContactMessage: "Want to talk it through? We're here",
		ContactEmail:   "hey@byrdhq.com",
	},
	GeneratedAt: time.Now(),
}

// renewalSuccessTemplate is a common template for renewal success emails
// Event: Triggered when a renewal payment succeeds
var renewalSucceededTemplate = CommonTemplate{
	PreviewText: "You just made our day (and yes, your renewal's confirmed)",
	Title:       "High fives all around",
	Subtitle:    "Your renewal just came through",
	Body: []string{
		"Hi there,",
		"Quick confession: we're supposed to send you a standard payment confirmation, but we're too thrilled for that. Your renewal just made our whole team's day.",
		"Beyond the usual business stuff, we wanted to say thanks for making Byrd better for everyone else. To show our appreciation, we've secretly snuck extra user seats - no charge (don't tell accounting).",
	},
	CTA: &CallToAction{
		ButtonText: "Let's Get Started",
		ButtonURL:  "https://byrdhq.com/dashboard",
		FooterText: "Here's to another year of crushing it together.",
	},
	Footer: Footer{
		ContactMessage: "Need help? We've got your back:",
		ContactEmail:   "hey@byrd.com",
	},
	GeneratedAt: time.Now(),
}

// renewalCanceledTemplate is a common template for renewal canceled emails
// Event: Triggered when a renewal payment is canceled
var renewalCanceledTemplate = CommonTemplate{
	PreviewText: "One last chat before we go",
	Title:       "We noticed you're wrapping things up with Byrd, and we wanted to say thanks",
	Body: []string{
		"Hi there,",
		"We noticed you're moving on from Byrd, and we respect that. But before you step away, we'd love to hear your thoughts on what we could have done better.",
		"Share your candid feedback - the hits and misses - and we'll give you a month's access on us. Maybe we can show you what's changed, or maybe we'll just learn how to do better next time.",
	},
	CTA: &CallToAction{
		ButtonText: "Share Feedback",
		ButtonURL:  "https://byrdhq.com/feedback",
		FooterText: "We're here to listen, and we're here to learn.",
	},
	Footer: Footer{
		ContactMessage: "Need a different arrangement? Lets talk:",
		ContactEmail:   "hey@byrdhq.com",
	},
	GeneratedAt: time.Now(),
}

// weeklyRoundupTemplate is a common template for weekly roundup emails
// Event: Triggered when a weekly roundup is sent
var weeklyRoundupTemplate = SectionedTemplate{
	Competitor:  "Competitor X",
	FromDate:    time.Now().AddDate(0, 0, -7),
	ToDate:      time.Now(),
	Summary:     "This week has been particularly active with major updates across branding and pricing. Here's what you need to know.",
	GeneratedAt: time.Now(),
	Sections: map[string]Section{
		"branding": {
			Title:   "BRANDING",
			Summary: "Major brand refresh and positioning updates",
			Bullets: []BulletPoint{
				{
					Text:    "Updated logo and visual identity",
					LinkURL: "#",
				},
				{
					Text:    "New brand guidelines released",
					LinkURL: "#",
				},
			},
		},
		"pricing": {
			Title:   "PRICING",
			Summary: "New pricing structure implemented",
			Bullets: []BulletPoint{
				{
					Text:    "Introduced new enterprise tier",
					LinkURL: "#",
				},
			},
		},
	},
}

var templates = map[TemplateName]Template{
	// -- user lifecycle templates ---
	WaitlistTemplate: &waitlistTemplate,
	WelcomeTemplate:  &welcomeTemplate,

	// -- trial lifecycle templates ---
	TrailSucceededTemplate: &trialSucceededTemplate,

	// -- workspace user lifecycle templates ---
	RequestWorkspaceInviteTemplate:   &requestWorkspaceInviteTemplate,
	PendingWorkspaceInviteTemplate:   &pendingWorkspaceInviteTemplate,
	AcceptWorkspaceInviteTemplateFTU: &acceptWorkspaceInviteTemplateFTU,
	AcceptWorkspaceInviteTemplateRU:  &acceptWorkspaceInviteTemplateRU,
	DeclineWorkspaceInviteTemplate:   &declineWorkspaceInviteTemplate,

	// -- renewal lifecycle templates ---
	RenewalFailedTemplate:    &renewalFailedTemplate,
	RenewalSucceededTemplate: &renewalSucceededTemplate,
	RenewalCanceledTemplate:  &renewalCanceledTemplate,

	// -- workspace lifecycle templates ---
	DeletedWorkspaceTemplate: &deletedWorkspaceTemplate,

	// -- weekly roundup template ---
	WeeklyRoundupTemplate: &weeklyRoundupTemplate,
}
