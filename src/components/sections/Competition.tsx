import ContentGrid from "../block/ContentGridBlock";
import ScreenshotBlock from "../block/ScreenshotBlock";
import SectionHeader from "../block/SectionHeaderBlock";

const CompetitionSection = () => {
	const gridContent = {
		mainTitle: {
			desktop: "There are no\nparticipation\ntrophies in\nbusiness",
			mobile: "There are no participation trophies in business",
		},
		columns: [
			{
				title: {
					desktop: "Your Product\nBetter Have Claws",
					mobile: "Your Product Better Have Claws",
				},
				description:
					"Chances are, the very thing you're working on, someone's already thought of it, tried it, or is working on it. The most dangerous competition is the one you pretend doesn't exist.",
				linkText: "What you don't know CAN hurt you",
				linkHref: "/",
			},
			{
				title: {
					desktop: "Second Place is\nFirst Loser",
					mobile: "Second Place is First Loser",
				},
				description:
					"First to market, last to profit? Your competitors are gunning for gold. Don't become a footnote in someone else's success story.",
				linkText: "Customers close, Competitors closer",
				linkHref: "/",
			},
			{
				title: {
					desktop: "Size isn't a Shield\nIt's a Target",
					mobile: "Size isn't a Shield It's a Target",
				},
				description:
					"Think you know your market? Think again. Being big doesn't make you invincible; it makes you visible. And in business, visibility without vigilance is a death sentence.",
				linkText: "Blindfolds Are for Piñatas, Not CEOs",
				linkHref: "/",
			},
		],
	};

	const headerContent = {
		title: "Your Company Has Competitors\nNot Just Cheerleaders",
		subtitle:
			"Ignorance isn't bliss. Transform competitive blind spots into your unfair advantage.",
	};

	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			{/* Hero Text Content */}
			<SectionHeader {...headerContent} />

			{/* Screenshot Container */}
			<ScreenshotBlock imageSrc="/inspector.png" />

			<ContentGrid {...gridContent} />
		</div>
	);
};

export default CompetitionSection;
