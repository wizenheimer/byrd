import ScreenshotBlock from "../block/ScreenshotBlock";
import SectionHeader from "../block/SectionHeaderBlock";
import SectionWithLead from "../block/SectionWithLead";

const ContextSection = () => {
	const headerContent = {
		title: "Your Company Has Competitors\nNot Just Cheerleaders",
		subtitle: "Turn Your Competitor's Moves Into Your Customer Wins",
	};

	const sectionContent = {
		leadText: {
			desktop: "There are No\nParticipation\nTrophies in\nBusiness",
		},
		contentColumns: [
			{
				title: {
					desktop: "Turn Competitors into Case Studies",
					mobile: "Turn Competitors into Case Studies",
				},
				description:
					"Every change, no matter how small, gets flagged. Whether it's roadmap changes, pricing shifts, or positioning plays. Consider yourself briefed.",
				linkText: "Competition Never Sleeps. But You Can",
				linkHref: "/onboarding?source=S2C1",
			},
			{
				title: {
					desktop: "Focus on Every Move That Matters",
					mobile: "Focus on Every Move That Matters",
				},
				description:
					"Your time is better spent crushing competitors, not cataloging them. Track what matters, skip what doesn't. Build upon real market context.",
				linkText: "Don't Play Catchup. Call the Shots",
				linkHref: "/onboarding?source=S2C2",
			},
		] as const,
	};
	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			{/* Hero Text Content */}
			<SectionHeader {...headerContent} />

			<ScreenshotBlock imageSrc="/competitive-intelligence.png" />

			{/* Bottom Three Column Section */}
			<SectionWithLead {...sectionContent} />
		</div>
	);
};

export default ContextSection;
