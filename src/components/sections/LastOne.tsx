import SectionHeader from "./block/SectionHeaderBlock";
import SplitBlock from "./block/SplitBlock";

const LastOne = () => {
	const headerContent = {
		title: "Don't be the last one to know",
		subtitle:
			"Spot emerging competitors long before they become existential threats.",
	};

	const splitContent = {
		leftColumn: {
			imageSrc: "/updates-mobile-view.png",
			title: "Shared competitive picture",
			paragraphs: [
				"Right now, you've got Marketing hoarding one set of competitor data, Sales clutching another, and Product off to something else. By the time this fragmented intel makes its way up the chain, competitors have already made their move, and you're playing on the back foot.",
				"With byrd product, marketing, sales, and even c-suite can all benefit from the same, constantly updated competitive picture.",
			],
			linkText: "Learn more about Unified Stream of Intel for Teams",
			linkHref: "/",
		},
		rightColumn: {
			imageSrc: "/integrations.png",
			title: "Enrich your knowledge base",
			paragraphs: [
				"Existing tools don't solve data silos; they create new ones and leave the plumbing for you. Byrd is built differently.",
				"Byrd offers first class integrations with the tools you already use - Slack, Notion, and even email. This means the competitive intelligence you gather doesn't just sit in yet another dashboard you'll forget to check. Instead, it flows seamlessly into your existing knowledge base and communication channels.",
			],
			linkText: "No new tabs, turn Slack into a war room",
			linkHref: "/",
		},
	};
	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			{/* Hero Text Content */}
			<SectionHeader {...headerContent} />
			<SplitBlock
				leftColumn={splitContent.leftColumn}
				rightColumn={splitContent.rightColumn}
			/>
		</div>
	);
};

export default LastOne;
