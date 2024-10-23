import React from "react";
import { Button } from "@/components/ui/button";
import PreviewContainer from "./PreviewContainer";

const ExitCTA = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-60">
			<div className="max-w-7xl mx-auto px-4 pt-24 pb-48">
				{/* Hero Text Content */}
				<div className="text-center max-w-3xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-6">
						We watch your competitors,
						<br />
						so you don't have to
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
							Is byrd for me?
						</Button>
					</div>
				</div>
			</div>
		</div>
	);
};

export default ExitCTA;
