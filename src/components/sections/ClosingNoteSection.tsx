import HeroSection from "../block/HeroBlock";

const ClosingNoteSection = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			<HeroSection
				title={{
					desktop: "Your Unfair Advantage Starts Here",
					mobile: "Your Unfair Advantage Starts Here",
				}}
				description="Your competitors are everywhere. So are we."
				primaryButton={{
					label: "Get Started",
					href: "/signup",
				}}
				secondaryButton={{
					label: "Contact Sales",
					href: "/contact",
				}}
			/>
		</div>
	);
};

export default ClosingNoteSection;
