// src/components/sections/ScaleupSection.tsx
import ScreenshotBlock from "../block/ScreenshotBlock";
import SectionHeader from "../block/SectionHeaderBlock";
import SectionWithLead from "../block/SectionWithLead";

const ScaleupSection = () => {
  const content = {
    leadText: {
      desktop: "The Market's Hungry,\nand You Look Like Lunch",
    },
    contentColumns: [
      {
        title: {
          desktop: "Differentiate Early.\nAnd Often.",
          mobile: "Differentiate Early. And Often.",
        },
        description:
          "Don't just think outside the box. Study every box out there and build a product that resonates best with your customers.",
        linkText: "Competitive Intelligence for Scaleups",
        linkHref: "/",
      },
      {
        title: {
          desktop: "Size Isn't a Shield.\nIt's a Target.",
          mobile: "Size Isn't a Shield. It's a Target.",
        },
        description:
          "Being big doesn't make you invincible; it makes you visible. And in business, visibility without vigilance is a death sentence.",
        linkText: "Disruption Doesn't Send Invites",
        linkHref: "/",
      },
    ] as const,
  };

  const headerContent = {
    title: "Differentiation Requires Context",
    subtitle:
      "You can't disrupt what you don't understand. Stop building with blindfolds on.",
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

export default ScaleupSection;
