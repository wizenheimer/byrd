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

			<div className="max-w-6xl mx-auto px-4 relative mt-20 md:mt-48">
				<div className="grid grid-cols-1 md:grid-cols-4 gap-8 md:gap-12">
					{/* Title Column - Hidden on mobile */}
					<div className="hidden md:block">
						<h2 className="text-3xl font-bold leading-tight">
							<span className="hidden md:inline">
								There are no
								<br />
								participation
								<br />
								trophies in
								<br />
								business
							</span>
							<span className="md:hidden">
								There are no participation trophies in business
							</span>
						</h2>
					</div>

					{/* Content Columns */}
					<div className="col-span-full md:col-span-3">
						<div className="grid grid-cols-1 md:grid-cols-3 gap-12 md:gap-12">
							{/* Column 1 */}
							<div className="space-y-6 md:grid md:grid-rows-[auto_1fr_auto] md:gap-8 max-w-sm mx-auto md:mx-0">
								<h2 className="text-3xl font-bold">
									<span className="hidden md:inline">
										Your Product
										<br />
										Better Have Claws
									</span>
									<span className="md:hidden">
										Your Product Better Have Claws
									</span>
								</h2>
								<p className="text-gray-600 leading-relaxed">
									Chances are, the very thing you&apos;re working on,
									someone&apos;s already thought of it, tried it, or is working
									on it. The most dangerous competition is the one you pretend
									doesn&apos;t exist.
								</p>
								<div>
									<Link
										href="/"
										className="text-gray-500 hover:text-gray-900 font-bold underline"
									>
										What you don&apos;t know CAN hurt you
									</Link>
								</div>
							</div>

							{/* Column 2 */}
							<div className="space-y-6 md:grid md:grid-rows-[auto_1fr_auto] md:gap-8 max-w-sm mx-auto md:mx-0">
								<h2 className="text-3xl font-bold">
									<span className="hidden md:inline">
										Second Place is
										<br />
										First Loser
									</span>
									<span className="md:hidden">Second Place is First Loser</span>
								</h2>
								<p className="text-gray-600 leading-relaxed">
									First to market, last to profit? Your competitors are gunning
									for gold. Don&apos;t become a footnote in someone else&apos;s
									success story.
								</p>
								<div>
									<Link
										href="/"
										className="text-gray-500 hover:text-gray-900 font-bold underline"
									>
										Customers close, Competitors closer
									</Link>
								</div>
							</div>

							{/* Column 3 */}
							<div className="space-y-6 md:grid md:grid-rows-[auto_1fr_auto] md:gap-8 max-w-sm mx-auto md:mx-0">
								<h2 className="text-3xl font-bold">
									<span className="hidden md:inline">
										Size isn&apos;t a Shield
										<br />
										It&apos;s a Target
									</span>
									<span className="md:hidden">
										Size isn&apos;t a Shield It&apos;s a Target
									</span>
								</h2>
								<p className="text-gray-600 leading-relaxed">
									Think you know your market? Think again. Being big
									doesn&apos;t make you invincible; it makes you visible. And in
									business, visibility without vigilance is a death sentence.
								</p>
								<div>
									<Link
										href="/"
										className="text-gray-500 hover:text-gray-900 font-bold underline"
									>
										Blindfolds Are for Piñatas, Not CEOs
									</Link>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default CompetitionSection;
