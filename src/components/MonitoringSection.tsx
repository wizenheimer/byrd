import Link from "next/link";
import React from "react";
import AnimatedImage from "./AnimatedImage";

const MonitoringSection = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-40">
			{/* Hero Text Content */}
			<div className="max-w-7xl mx-auto px-4 pt-20 pb-40">
				<div className="text-center max-w-4xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-6">
						Your Company Has Competitors
						<br />
						Not Just Cheerleaders
					</h1>
					<p className="text-lg text-gray-600 mt-8">
						Ignorance isn&apos;t bliss. Transform competitive blind spots into
						your unfair advantage.
					</p>
				</div>
			</div>

			{/* Screenshot Container */}
			<div className="max-w-6xl mx-auto px-4 relative ">
				<AnimatedImage imageSrc="/web-monitoring.png" />
			</div>

			{/* Bottom Three Column Section */}
			<div className="max-w-7xl mx-auto px-4 relative mt-20 md:mt-48">
				<div className="grid grid-cols-1 md:grid-cols-3 gap-10 md:gap-20">
					{/* Column 1 */}
					<div className="hidden md:block md:ml-20">
						<h2 className="text-3xl font-semibold">
							Stop admiring
							<br />
							the problem
							<br />
							and start
							<br />
							solving it.
						</h2>
					</div>

					{/* Column 2 and 3 wrapper */}
					<div className="md:col-span-2 grid grid-cols-1 sm:grid-cols-2 gap-10 md:gap-20">
						{/* Column 2 */}
						<div className="max-w-sm mx-auto sm:mx-0">
							<h2 className="text-2xl font-semibold mb-4">
								Your Rivals Don&apos;t Rest
							</h2>
							<p className="text-base font-medium">
								You think you&apos;re crushing it? So does everyone else. Stop
								patting yourself on the back and start learning from
								competitors. Build a product that resonates best with your
								customers.
							</p>
							<div className="mt-8 md:mt-12">
								<Link
									href="/"
									className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
								>
									Learn more about competitive intelligence
								</Link>
							</div>
						</div>

						{/* Column 3 */}
						<div className="max-w-sm mx-auto sm:mx-0">
							<h2 className="text-2xl font-semibold mb-4">
								Be Competitor Informed
							</h2>
							<p className="text-base font-medium">
								Your competitors are constantly evolving, changing prices,
								launching new features. It&apos;s not what you know about your
								competitors - it&apos;s what you don&apos;t know yet that
								matters.
							</p>
							<div className="mt-8 md:mt-12">
								<Link
									href="/"
									className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
								>
									Learn more about market intelligence
								</Link>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default MonitoringSection;
