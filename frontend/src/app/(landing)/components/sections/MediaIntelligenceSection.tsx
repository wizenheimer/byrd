// src/components/sections/MediaIntelligenceSection.tsx
import ContentGrid from "../block/ContentGridBlock";
import ScreenshotBlock from "../block/ScreenshotBlock";
import SectionHeader from "../block/SectionHeaderBlock";

const MediaIntelligenceSection = () => {
	const headerContent = {
		title: "Track the Stories They Break\nAnd Those They Don't",
		subtitle:
			"Your front-row seat to their media strategy.\nFrom social chatter to media mentions, catch every signal that matters.",
	};

	const gridContent = {
		mainTitle: {
			desktop: "Observe Every Competitor Move Before Your Customers Do",
			mobile: "Observe Every Competitor Move Before Your Customers Do",
		},
		columns: [
			{
				title: {
					desktop: "Markets Talk.\nWe Translate.",
					mobile: "Markets Talk. We Translate.",
				},
				description:
					"Measure what matters. See where they publish. Track what resonates. Learn why it works.",
				linkText: "Track Their Narrative Before It Hits Mainstream",
				linkHref: "/",
			},
			{
				title: {
					desktop: "First To Know.\nFirst To Move.",
					mobile: "First To Know. First To Move.",
				},
				description:
					"Track every dime or dollar your competitors raise. Keep a pulse on leadership changes, partnerships, and more.",
				linkText: "Start Pitching with Precision, not Prayers",
				linkHref: "/",
			},
			{
				title: {
					desktop: "Guess Less.\nLearn More.",
					mobile: "Guess Less. Learn More.",
				},
				description:
					"From major announcements to minor mentions. Never miss on coverage that shapes market perceptions or influences buyers.",
				linkText: "Market Visibility Without the Busy Work",
				linkHref: "/",
			},
		],
	};

	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			<SectionHeader {...headerContent} />

			<ScreenshotBlock imageSrc="/media-intelligence.png" />

			<ContentGrid {...gridContent} />
		</div>
	);
};

export default MediaIntelligenceSection;
