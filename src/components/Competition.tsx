import Link from "next/link";
import React from "react";
import AnimatedImage from "./AnimatedImage";

const CompetitionSection = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-40">
			{/* Hero Text Content */}
			<div className="max-w-7xl mx-auto px-4 pt-20 pb-40">
				<div className="text-center max-w-4xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-6">
						Competition Who Cares?
					</h1>
					<p className="text-lg text-gray-600 mt-8">
						Your customers do. They compare. They seek alternatives.
					</p>
				</div>
			</div>

			{/* Screenshot Container */}
			<div className="max-w-6xl mx-auto px-4 relative">
				<AnimatedImage imageSrc="/inspector.png" />
			</div>

			<div className="max-w-6xl mx-auto px-4 relative mt-48">
				<div className="grid grid-cols-4 gap-12">
					{/* Title Column */}
					<div>
						<h2 className="text-3xl font-bold leading-tight">
							There are no
							<br />
							participation
							<br />
							trophies in
							<br />
							business
						</h2>
					</div>

					{/* Content Columns */}
					<div className="col-span-3 grid grid-cols-3 gap-12">
						{/* Column 1 */}
						<div className="flex flex-col">
							<div>
								<h3 className="text-xl font-semibold mb-4">
									Your Product
									<br />
									Better Have Claws
								</h3>
								<p className="text-gray-600 mb-6">
									Chances are, the very thing you&apos;re working on,
									someone&apos;s already thought of it, tried it, or is working
									on it. The most dangerous competition is the one you pretend
									doesn&apos;t exist.
								</p>
							</div>
							<div className="mt-auto text-sm">
								<Link
									href="/"
									className="text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
								>
									What you don&apos;t know CAN hurt you
								</Link>
							</div>
						</div>

						{/* Column 2 */}
						<div className="flex flex-col">
							<div>
								<h3 className="text-xl font-semibold mb-4">
									Second Place is
									<br />
									First Loser
								</h3>
								<p className="text-gray-600 mb-6">
									First to market, last to profit? Your competitors are gunning
									for gold. Don&apos;t become a footnote in someone else&apos;s
									success story.
								</p>
							</div>
							<div className="mt-auto text-sm">
								<Link
									href="/"
									className="text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
								>
									Customers close, Competitors closer
								</Link>
							</div>
						</div>

						{/* Column 3 */}
						<div className="flex flex-col">
							<div>
								<h3 className="text-xl font-semibold mb-4">
									Size isn&apos;t a Shield
									<br />
									It&apos;s a Target
								</h3>
								<p className="text-gray-600 mb-6">
									Think you know your market? Think again. Being big
									doesn&apos;t make you invincible; it makes you visible. And in
									business, visibility without vigilance is a death sentence.
								</p>
							</div>
							<div className="mt-auto text-sm">
								<Link
									href="/"
									className="text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
								>
									Blindfolds Are for Piñatas, Not CEOs
								</Link>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default CompetitionSection;
