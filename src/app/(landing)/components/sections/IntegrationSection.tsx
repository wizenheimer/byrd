// src/components/sections/IntegrationSection.tsx
import SectionHeader from "../block/SectionHeaderBlock";
import SplitBlock from "../block/SplitBlock";

const IntegrationSection = () => {
  const headerContent = {
    title: "Don't Be the Last One To Know",
    subtitle:
      "Spot emerging competitors long before they become existential threats.\nNo more of 'how did we miss that?' moments.",
  };

  const splitContent = {
    leftColumn: {
      imageSrc: "/inbox-mobile-view.png",
      title: "Winning Is a Team Sport",
      paragraphs: [
        "Bring in your heavy hitters. Keep everyone on the same page. When everyone sees the same signals, nobody misses the big moves.",
        "Deliver real-time updates across all your teams, with customized views for each role. No more of 'marketing said this, sales heard that'.",
      ],
      linkText: "Stop losing intel in endless email threads",
      linkHref: "/",
    },
    rightColumn: {
      imageSrc: "/integrations.png",
      title: "Same Tools. No New Dashboards.",
      paragraphs: [
        "Existing tools don't solve data silos; they create new ones and leave the plumbing for you. Byrd is built differently.",
        "With first class integrations, competitive intelligence you gather no longer sits in yet another dashboard you'll forget to check.",
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

export default IntegrationSection;
