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
					label: "Get started",
					href: "/onboarding",
				}}
				secondaryButton={{
					label: "Talk to a human",
					href: "https://cal.com/nayann/byrd",
				}}
			/>

			<ScreenshotBlock imageSrc="/overview.png" />
		</div>
	);
};

export default Hero;
