import ScreenshotBlock from "./block/ScreenshotBlock";
import SectionHeader from "./block/SectionHeaderBlock";
import SectionWithLead from "./block/SectionWithLead";

const MonitoringSection = () => {
	const headerContent = {
		title: "Your Company Has Competitors\nNot Just Cheerleaders",
		subtitle:
			"Ignorance isn't bliss. Transform competitive blind spots into your unfair advantage.",
	};

	const sectionContent = {
		leadText: {
			desktop: "Stop admiring\nthe problem\nand start\nsolving it.",
		},
		contentColumns: [
			{
				title: {
					desktop: "Your Rivals Don't Rest",
					mobile: "Your Rivals Don't Rest",
				},
				description:
					"You think you're crushing it? So does everyone else. Stop patting yourself on the back and start learning from competitors. Build a product that resonates best with your customers.",
				linkText: "Competitive Intelligence for Scaleups",
				linkHref: "/",
			},
			{
				title: {
					desktop: "Be Competitor Informed.",
					mobile: "Be Competitor Informed.",
				},
				description:
					"Your competitors are constantly evolving, changing prices, launching new features. It's not what you know about your competitors - it's what you don't know yet that matters.",
				linkText: "Don't play catchup. Call the shots.",
				linkHref: "/",
			},
		] as const,
	};
	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			{/* Hero Text Content */}
			<SectionHeader {...headerContent} />

			<ScreenshotBlock imageSrc="/web-monitoring.png" />

			{/* Bottom Three Column Section */}
			<SectionWithLead {...sectionContent} />
		</div>
	);
};

export default MonitoringSection;
