import ContentGrid from "../block/ContentGridBlock";
import ScreenshotBlock from "../block/ScreenshotBlock";
import SectionHeader from "../block/SectionHeaderBlock";

const CustomerIntelligenceSection = () => {
	const headerContent = {
		title: "Listen to the Voices\nThey're Ignoring",
		subtitle:
			"Every complaint is a competitive opportunity.\nTurn their user feedback into your feature wins.",
	};

	const gridContent = {
		mainTitle: {
			desktop: "Stop admiring\nthe problem and\nstart solving it",
			mobile: "Stop admiring the problem and start solving it",
		},
		columns: [
			{
				title: {
					desktop: "From Rage To Riches.",
					mobile: "From Rage To Riches.",
				},
				description:
					"Spot product gaps before they fix them. Learn what users love (and hate) about your competitors. Turn reviews into revenue.",
				linkText: "Their Users Are Talking. Are You Listening ?",
				linkHref: "/onboarding?source=S6C1",
			},
			{
				title: {
					desktop: "They Churn. You Learn.",
					mobile: "They Churn. You Learn.",
				},
				description:
					"Turn their dissatisfied users into your would be customers. While they're busy defending 1-star reviews, build what their customers wish they had.",
				linkText: "Let their Customers Inform your Roadmap",
				linkHref: "/onboarding?source=S6C2",
			},
			{
				title: {
					desktop: "Their Logos. Your Leads.",
					mobile: "Their Logos. Your Leads.",
				},
				description:
					"Stop prospecting cold. Every new case study your partners release is a fresh target for your pipeline. Turn their logos into your leads.",
				linkText: "Rivals are leaving money on the table. Go cash in",
				linkHref: "/onboarding?source=S6C3",
			},
		],
	};

	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			<SectionHeader {...headerContent} />

			<ScreenshotBlock imageSrc="/customer-intelligence.png" />

			<ContentGrid {...gridContent} />
		</div>
	);
};

export default CustomerIntelligenceSection;
