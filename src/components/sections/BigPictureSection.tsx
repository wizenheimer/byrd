import FeatureGrid from "../block/FeatureGridBlock";
import ScreenshotBlock from "../block/ScreenshotBlock";
import SectionHeader from "../block/SectionHeaderBlock";

const BigPictureSection = () => {
	const headerContent = {
		title: "Focus on the big picture\nLet us handle the details",
		subtitle:
			"Every hour spent tracking competitors is an hour spent not crushing them.\nStop wasting your brilliance on spreadsheets.",
	};

	const featureContent = {
		columns: [
			{
				title: {
					desktop: "More Bodies\nMore Problems",
					mobile: "More Bodies. More Problems",
				},
				description:
					"Planning on hiring an army of interns? Flooding them with busywork? Measuring outcome by the pounds of reports generated? Forget PMF, you just re-discovered the express lane to Chapter 11.",
				link: {
					text: {
						desktop: "Don't solve problems\nby creating new ones",
						mobile: "Don't solve problems, by creating new ones",
					},
					href: "/",
				},
			},
			{
				title: {
					desktop: "All Signal\nNo Noise",
					mobile: "All Signal. No Noise",
				},
				description:
					"With Byrd, each alert is like a personal briefing. Whether it's product updates, branding overhauls, or partnership plays, you'll be the first to know.",
				link: {
					text: {
						desktop: "Competition never sleeps,\nbut you can",
						mobile: "Competition never sleeps, but you can",
					},
					href: "/",
				},
			},
			{
				title: {
					desktop: "Real-time\nAll the Time",
					mobile: "Real-time. All the Time",
				},
				description:
					"With byrd you get clear, concise, and actionable intel delivered straight to you. It's like having a team of consultants on call, minus the expensive suits and PowerPoint.",
				link: {
					text: {
						desktop: "The Market's Hungry,\nand You Look Like Lunch",
						mobile: "The Market's Hungry, and You Look Like Lunch",
					},
					href: "/",
				},
			},
		],
	};

	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			{/* Hero Text Content */}
			<SectionHeader {...headerContent} />

			{/* Screenshot Container */}
			<ScreenshotBlock imageSrc="/newsletter-monitoring.png" />

			{/* Bottom Section - Aligned with PreviewContainer */}
			<FeatureGrid columns={featureContent.columns} />
		</div>
	);
};

export default BigPictureSection;
