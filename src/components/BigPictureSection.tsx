import Link from "next/link";
import React from "react";
import AnimatedImage from "./AnimatedImage";

const BigPictureSection = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-40">
			{/* Hero Text Content */}
			<div className="max-w-7xl mx-auto px-4 pt-20 pb-40">
				<div className="text-center max-w-4xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-6">
						Focus on the big picture
						<br />
						Let us handle the details
					</h1>
					<p className="text-lg text-gray-600 mt-8">
						Every hour spent tracking competitors is an hour spent not crushing
						them.
						<br />
						Stop wasting your brilliance on spreadsheets.
					</p>
				</div>
			</div>

			{/* Screenshot Container */}
			<div className="max-w-6xl mx-auto px-4 relative ">
				<AnimatedImage imageSrc="/newsletter-monitoring.png" />
			</div>

			{/* Bottom Section - Aligned with PreviewContainer */}
			<div className="max-w-6xl mx-auto px-4 relative mt-16 md:mt-32">
				<div className="grid grid-cols-1 md:grid-cols-3 gap-12 md:gap-24">
					{/* Column 1 */}
					<div className="space-y-6 md:grid md:grid-rows-[auto_1fr_auto] md:gap-8 max-w-sm mx-auto md:mx-0">
						<h2 className="text-2xl font-bold text-center md:text-left">
							<span className="hidden md:inline">
								More Bodies
								<br />
								More Problems
							</span>
							<span className="md:hidden">More Bodies More Problems</span>
						</h2>
						<p className="text-gray-600 leading-relaxed">
							Planning on hiring an army of interns? Flooding them with
							busywork? Measuring outcome by the pounds of reports generated?
							Forget PMF, you just re-discovered the express lane to Chapter 11.
						</p>
						<div className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700">
							<Link href="/" className="hidden md:inline">
								Don&apos;t solve problems
								<br />
								by creating new ones
							</Link>
							<Link href="/" className="md:hidden">
								Don&apos;t solve problems, by creating new ones
							</Link>
						</div>
					</div>

					{/* Column 2 */}
					<div className="space-y-6 md:grid md:grid-rows-[auto_1fr_auto] md:gap-8 max-w-sm mx-auto md:mx-0">
						<h2 className="text-2xl font-bold text-center md:text-left">
							<span className="hidden md:inline">
								All Signal
								<br />
								No Noise
							</span>
							<span className="md:hidden">All Signal No Noise</span>
						</h2>
						<p className="text-gray-600 leading-relaxed">
							With Byrd, each alert is like a personal briefing. Whether
							it&apos;s product updates, branding overhauls, or partnership
							plays, you&apos;ll be the first to know.
						</p>
						<div className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700">
							<Link href="/" className="hidden md:inline">
								Competition never sleeps,
								<br />
								but you can
							</Link>
							<Link href="/" className="md:hidden">
								Competition never sleeps, but you can
							</Link>
						</div>
					</div>

					{/* Column 3 */}
					<div className="space-y-6 md:grid md:grid-rows-[auto_1fr_auto] md:gap-8 max-w-sm mx-auto md:mx-0">
						<h2 className="text-2xl font-bold text-center md:text-left">
							<span className="hidden md:inline">
								Real-time
								<br />
								All the Time
							</span>
							<span className="md:hidden">Real-time All the Time</span>
						</h2>
						<p className="text-gray-600 leading-relaxed">
							With byrd you get clear, concise, and actionable intel delivered
							straight to you. It&apos;s like having a team of consultants on
							call, minus the expensive suits and PowerPoint.
						</p>
						<div className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700">
							<Link href="/" className="hidden md:inline">
								The Market&apos;s Hungry,
								<br />
								and You Look Like Lunch
							</Link>
							<Link href="/" className="md:hidden">
								The Market&apos;s Hungry, and You Look Like Lunch
							</Link>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default BigPictureSection;
