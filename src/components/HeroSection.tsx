import { Button } from "@/components/ui/button";
import React from "react";
import AnimatedImage from "./AnimatedImage";

const Hero = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden">
			<div className="max-w-7xl mx-auto px-4 pt-24 pb-16 md:pb-32 lg:pb-48">
				{/* Beta Badge */}
				<div className="flex justify-center mb-8">
					<div className="bg-black/10 rounded-full px-3 py-1 inline-flex items-center gap-2">
						<div className="w-2 h-2 rounded-full bg-green-500 animate-[statusBlink_3s_ease-in-out_infinite]" />
						<span className="text-sm">Beta</span>
					</div>
				</div>

				{/* Hero Text Content */}
				<div className="text-center max-w-3xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-6">
						We Watch Your Competitors,
						<br />
						So You Don&apos;t Have To
					</h1>
					<p className="text-lg text-gray-600 mb-8">
						Know the moves your competitors make, long before their own
						employees do.
					</p>
					<div className="flex gap-4 justify-center">
						<Button className="bg-black text-white hover:bg-black/90 px-8">
							Get Started
						</Button>
						<Button variant="outline" className="border-gray-200">
							Contact Sales
						</Button>
					</div>
				</div>
			</div>

			{/* Screenshot Container */}
			<div className="max-w-6xl mx-auto px-4 relative">
				<AnimatedImage imageSrc="/overview.png" />
			</div>
		</div>
	);
};

export default Hero;
