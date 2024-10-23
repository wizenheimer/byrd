import React from "react";
import PreviewContainer from "./PreviewContainer";
import Link from "next/link";

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
				<PreviewContainer
					imageSrc="/assets/dashboard.png"
					caption="Observe market shifts as they happen."
				/>
			</div>

			{/* Bottom Section - Aligned with PreviewContainer */}
			<div className="max-w-6xl mx-auto px-4 relative mt-32">
				<div className="grid grid-cols-3 gap-24">
					{/* Column 1 */}
					<div className="grid grid-rows-[auto_1fr_auto] gap-8">
						<h2 className="text-3xl font-semibold">
							More Bodies
							<br />
							More Problems
						</h2>
						<p className="text-base text-gray-600 leading-relaxed">
							Planning on hiring an army of interns? Flooding them with
							busywork? Measuring outcome by the pounds of reports generated?
							Forget PMF, you just re-discovered the express lane to Chapter 11.
						</p>
						<div className="inline-block text-gray-500 hover:text-gray-900 font-medium underline">
							<Link href="/">Learn more about market intelligence</Link>
						</div>
					</div>

					{/* Column 2 */}
					<div className="grid grid-rows-[auto_1fr_auto] gap-8">
						<h2 className="text-3xl font-semibold">
							All Signal
							<br />
							No Noise
						</h2>
						<p className="text-base text-gray-600 leading-relaxed">
							With Byrd, each alert is like a personal briefing. Whether it's
							product updates, branding overhauls, or partnership plays, you'll
							be the first to know.
						</p>
						<div className="inline-block text-gray-500 hover:text-gray-900 font-medium underline">
							<Link href="/">Learn more about market intelligence</Link>
						</div>
					</div>

					{/* Column 3 */}
					<div className="grid grid-rows-[auto_1fr_auto] gap-8">
						<h2 className="text-3xl font-semibold">
							Real-time
							<br />
							All the Time
						</h2>
						<p className="text-base text-gray-600 leading-relaxed">
							With byrd you get clear, concise, and actionable intel delivered
							straight to you. It's like having a team of consultants on call,
							minus the expensive suits and PowerPoint.
						</p>
						<div className="inline-block text-gray-500 hover:text-gray-900 font-medium underline">
							<Link href="/">Learn more about market intelligence</Link>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default BigPictureSection;
