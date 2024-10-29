import HeroSection from "../block/HeroBlock";

const ExitCTA = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			<HeroSection
				title={{
					desktop: "Your Unfair Advantage Starts Here",
					mobile: "Your Unfair Advantage Starts Here",
				}}
				description="Your rivals are making moves. Time to put them under surveillance."
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

export default ExitCTA;
