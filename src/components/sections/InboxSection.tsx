import ContentGrid from "../block/ContentGridBlock";
import ScreenshotBlock from "../block/ScreenshotBlock";
import SectionHeader from "../block/SectionHeaderBlock";

const InboxSection = () => {
	const headerContent = {
		title: "Monitor the Direct Line\nto Their Customers",
		subtitle:
			"Your competitors' best marketing strategies live in their email campaigns.\nLearn what they're telling their best customers.",
	};

	const gridContent = {
		mainTitle: {
			desktop: "Watch Their Messaging Evolve, Email by Email",
			mobile: "Watch Their Messaging Evolve, Email by Email",
		},
		columns: [
			{
				title: {
					desktop: "Never Miss a Message",
					mobile: "Never Miss a Message",
				},
				description:
					"Every competitor email archived and searchable, from day one. Search across all competitor emails by topic, product, or campaign.",
				linkText: "Learn Which Bets Your Competitors Are Making Next",
				linkHref: "/",
			},
			{
				title: {
					desktop: "Know When They Strike",
					mobile: "Know When They Strike",
				},
				description:
					"Track exactly when they hit send. Learn which days and times their emails land hardest. Beat them to the inbox next time.",
				linkText: "Focus on What Works. Skip the Rest",
				linkHref: "/",
			},
			{
				title: {
					desktop: "Bookmark their winning moves",
					mobile: "Bookmark their winning moves",
				},
				description:
					"Build Your Swipe Files. Turn their best work into your advantage. Give your marketing team the ammo they need.",
				linkText: "Your Competition's Greatest Hits. Save Now. Strike Later",
				linkHref: "/",
			},
		],
	};

	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			<SectionHeader {...headerContent} />

			<ScreenshotBlock imageSrc="/inbox-monitoring.png" />

			<ContentGrid {...gridContent} />
		</div>
	);
};

export default InboxSection;
