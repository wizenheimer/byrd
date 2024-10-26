import Link from "next/link";
import React from "react";
import AnimatedImage from "./AnimatedImage";

const DifferentiationSection = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-32 md:mt-48 lg:mt-60">
			{/* Hero Text Content */}
			<div className="max-w-7xl mx-auto px-4 pt-12 md:pt-16 lg:pt-20 pb-20 md:pb-32 lg:pb-40">
				<div className="text-center max-w-4xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-6">
						Differentiation Requires Context
					</h1>
					<p className="text-lg text-gray-600 mt-8">
						You can&apos;t disrupt a market you don&apos;t understand, and you
						can&apos;t understand a market without knowing who the players are
						and what they&apos;re doing.
					</p>
				</div>
			</div>

			{/* Screenshot Container */}
			<div className="max-w-6xl mx-auto px-4 relative">
				<AnimatedImage imageSrc="/differentiation.png" />
			</div>

			{/* Bottom Three Column Section - Aligned with PreviewContainer */}
			<div className="max-w-6xl mx-auto px-4 relative mt-12 md:mt-32 lg:mt-48">
				<div className="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-20">
					{/* Column 1 */}
					<div className="md:block grid grid-rows-[1fr] items-start">
						<h2 className="text-3xl font-semibold">
							<span className="hidden md:inline">
								There&apos;s nothing
								<br />
								like being too
								<br />
								early, there&apos;s
								<br />
								only too late.
							</span>
						</h2>
					</div>

					{/* Column 2 */}
					<div className="grid grid-rows-[auto_1fr_auto] gap-6 md:gap-8 max-w-sm mx-auto md:mx-0">
						<h2 className="text-2xl font-semibold">
							<span className="hidden md:inline">
								Know your Goliath
								<br />
								inside out.
							</span>
							<span className="md:hidden">Know your Goliath inside out.</span>
						</h2>
						<p className="text-base font-medium text-gray-600 leading-relaxed">
							You think you&apos;re crushing it? So does everyone else. Stop
							patting yourself on the back and start learning from competitors.
							Build a product that resonates best with your customers.
						</p>
						<div>
							<Link
								href="/"
								className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
							>
								Competitive Intelligence for Scaleups
							</Link>
						</div>
					</div>

					{/* Column 3 */}
					<div className="grid grid-rows-[auto_1fr_auto] gap-6 md:gap-8 max-w-sm mx-auto md:mx-0">
						<h2 className="text-2xl font-semibold">
							<span className="hidden md:inline">
								Differentiate Early.
								<br />
								And Often.
							</span>
							<span className="md:hidden">Differentiate Early. And Often.</span>
						</h2>
						<p className="text-base font-medium text-gray-600 leading-relaxed">
							Your competitors are constantly evolving, changing prices,
							launching new features. It&apos;s not what you know about your
							competitors - it&apos;s what you don&apos;t know yet that matters.
						</p>
						<div>
							<Link
								href="/"
								className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
							>
								Don&apos;t play catchup. Call the shots.
							</Link>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default DifferentiationSection;
