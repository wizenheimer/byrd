import ScreenshotBlock from "./block/ScreenshotBlock";
import SectionHeader from "./block/SectionHeaderBlock";
import SectionWithLead from "./block/SectionWithLead";

const DifferentiationSection = () => {
	const content = {
		leadText: {
			desktop:
				"There's nothing\nlike being too\nearly, there's\nonly too late.",
		},
		contentColumns: [
			{
				title: {
					desktop: "Know your Goliath\ninside out.",
					mobile: "Know your Goliath inside out.",
				},
				description:
					"You think you're crushing it? So does everyone else. Stop patting yourself on the back and start learning from competitors. Build a product that resonates best with your customers.",
				linkText: "Competitive Intelligence for Scaleups",
				linkHref: "/",
			},
			{
				title: {
					desktop: "Differentiate Early.\nAnd Often.",
					mobile: "Differentiate Early. And Often.",
				},
				description:
					"Your competitors are constantly evolving, changing prices, launching new features. It's not what you know about your competitors - it's what you don't know yet that matters.",
				linkText: "Don't play catchup. Call the shots.",
				linkHref: "/",
			},
		] as const,
	};

	const headerContent = {
		title: "Differentiation Requires Context",
		subtitle:
			"You can't disrupt a market you don't understand, and you can't understand a market without knowing who the players are and what they're doing.",
	};

	return (
		<div className="w-full bg-background relative overflow-hidden mt-32 md:mt-48 lg:mt-60">
			{/* Hero Text Content */}
			<SectionHeader {...headerContent} />

			{/* Screenshot Container */}
			<ScreenshotBlock imageSrc="/differentiation.png" />

			{/* Bottom Three Column Section - Aligned with PreviewContainer */}
			<SectionWithLead {...content} />
		</div>
	);
};

export default DifferentiationSection;
