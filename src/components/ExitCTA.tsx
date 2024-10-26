import { Button } from "@/components/ui/button";
import React from "react";

const ExitCTA = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-20 md:mt-32 lg:mt-40">
			<div className="max-w-7xl mx-auto px-4 pt-16 md:pt-20 lg:pt-24 pb-20 md:pb-32 lg:pb-40">
				{/* Hero Text Content */}
				<div className="text-center max-w-3xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-6">
						Build Your Product on High Ground
					</h1>
					<p className="text-lg text-gray-600 mb-8">
						No noise, no fluff, just the intel you need to stay ahead.
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
		</div>
	);
};

export default ExitCTA;
