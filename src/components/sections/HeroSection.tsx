import HeroSection from "../block/HeroBlock";
import ScreenshotBlock from "../block/ScreenshotBlock";

const Hero = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden">
			<HeroSection
				title={{
					desktop: "We Watch Your Competitors,\nSo You Don't Have To",
					mobile: "We Watch Your Competitors, So You Don't Have To",
				}}
				description="Know the moves your competitors make, long before their own employees do"
				badge={{
					text: "Beta",
				}}
				primaryButton={{
					label: "Get Started",
					href: "/get-started",
				}}
				secondaryButton={{
					label: "Contact Sales",
          href: "https://cal.com/nayann/byrd",
				}}
			/>

			<ScreenshotBlock imageSrc="/overview.png" />
		</div>
	);
};

export default Hero;
