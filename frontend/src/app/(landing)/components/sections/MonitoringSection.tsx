import ContentGrid from "../block/ContentGridBlock";
import ScreenshotBlock from "../block/ScreenshotBlock";
import SectionHeader from "../block/SectionHeaderBlock";

const MonitoringSection = () => {
	const headerContent = {
		title: "Focus on the Big Picture\nLet Us Handle the Details",
		subtitle:
			"Every hour spent tracking competitors is an hour spent not crushing them.\nStop wasting your brilliance on spreadsheets.",
	};

	const gridContent = {
		mainTitle: {
			desktop: "Your Competitors'\nNightmares\n Just Got Real",
			mobile: "Your Competitors' Nightmares Just Got Real",
		},
		columns: [
			{
				title: {
					desktop: "All Signal.\nNo Noise.",
					mobile: "All Signal. No Noise.",
				},
				description:
					"Mine competitor mistakes for your advantage. Skip the guesswork. Start with what already works.",
				linkText: "Keep Your Customers Close, and Competitors Closer",
				linkHref: "/onboarding?source=S3C1",
			},
			{
				title: {
					desktop: "Real-Time.\nAll the Time.",
					mobile: "Real-Time. All the Time.",
				},
				description:
					"Make alerts your allies. Stay ahead of every product move. Stop discovering competitor announcements through lost deals.",
				linkText: "Stop Losing Deals to Competitors' Price Changes",
				linkHref: "/onboarding?source=S3C2",
			},
			{
				title: {
					desktop: "Their Next Move.\nYour Next Win.",
					mobile: "Their Next Move. Your Next Win.",
				},
				description:
					"Skip watching from the sidelines. Head straight to what works. Start turning their moves into your strategic wins.",
				linkText: "Beat them at their own game, every single time",
				linkHref: "/onboarding?source=S3C3",
			},
		],
	};

	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			<SectionHeader {...headerContent} />

			<ScreenshotBlock imageSrc="/product-intelligence.png" />

			<ContentGrid {...gridContent} />
		</div>
	);
};

export default MonitoringSection;
