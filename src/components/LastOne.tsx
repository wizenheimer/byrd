import React from "react";
import PreviewContainer from "./PreviewContainer";
import Link from "next/link";

const LastOne = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-40">
			{/* Hero Text Content */}
			<div className="max-w-7xl mx-auto px-4 pt-20 pb-20">
				<div className="text-center max-w-4xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-4">
						Don't be the last one to know.
					</h1>
					<p className="text-lg text-gray-600 mt-8">
						Spot emerging competitors long before they become existential
						threats.
						<br />
						No more of 'How did we miss that?' moments.
					</p>
				</div>
			</div>

			<div className="max-w-7xl mx-auto">
				{/* Two Column Layout with PreviewContainers */}
				<div className="grid grid-cols-1 lg:grid-cols-2 gap-24 px-8 py-12">
					{/* Left Column */}
					<div>
						<div className="relative mb-16">
							<PreviewContainer
								imageSrc="/assets/dashboard.png"
								variant="fitted"
								parallaxStrength={0}
							/>
						</div>
						<div className="space-y-6">
							<h2 className="text-3xl font-bold">Shared competitive picture</h2>
							<p className="text-gray-600">
								Right now, you've got Marketing hoarding one set of competitor
								data, Sales clutching another, and Product off to something
								else. By the time this fragmented intel makes its way up the
								chain, competitors have already made their move, and you're
								playing on the back foot.
							</p>
							<p className="text-gray-600">
								With byrd product, marketing, sales, and even c-suite can all
								benefit from the same, constantly updated competitive picture.
							</p>
						</div>
						<div className="mt-12">
							<Link
								href="/"
								className="text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
							>
								Learn more about Unified Stream of Intel for Teams
							</Link>
						</div>
					</div>

					{/* Right Column */}
					<div>
						<div className="relative mb-16">
							<PreviewContainer
								imageSrc="/assets/dashboard.png"
								variant="fitted"
								parallaxStrength={0}
							/>
						</div>
						<div className="space-y-6">
							<h2 className="text-3xl font-bold">Enrich your knowledge base</h2>
							<p className="text-gray-600">
								Existing tools don't solve data silos; they create new ones and
								leave the plumbing for you. Byrd is built differently.
							</p>
							<p className="text-gray-600">
								Byrd offers first class integrations with the tools you already
								use - Slack, Notion, and even email. This means the competitive
								intelligence you gather doesn't just sit in yet another
								dashboard you'll forget to check. Instead, it flows seamlessly
								into your existing knowledge base and communication channels.
							</p>
						</div>
						<div className="mt-12">
							<Link
								href="/"
								className="text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
							>
								No new tabs, turn Slack into a war room
							</Link>
						</div>
					</div>
				</div>
			</div>
			{/* Bottom Three Column Section - Aligned with PreviewContainer */}
			<div className="max-w-6xl mx-auto px-4 relative mt-48">
				{" "}
				{/* Changed from max-w-7xl to max-w-6xl */}
				<div className="grid grid-cols-3 gap-20">
					{/* Column 1 */}
					<div className="grid grid-rows-[1fr] items-start">
						{" "}
						{/* Removed ml-20 */}
						<h2 className="text-3xl font-semibold">
							There's nothing like being too early,
							<br />
							there's only too late.
						</h2>
					</div>

					{/* Column 2 */}
					<div className="grid grid-rows-[auto_1fr_auto] gap-8">
						{" "}
						{/* Removed max-w-sm */}
						<h2 className="text-2xl font-semibold">
							Know your Goliath inside out.
						</h2>
						<p className="text-base font-medium text-gray-600 leading-relaxed">
							You think you're crushing it? So does everyone else. Stop patting
							yourself on the back and start learning from competitors. Build a
							product that resonates best with your customers.
						</p>
						<div>
							<Link
								href="/"
								className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
							>
								Learn more about competitive intelligence
							</Link>
						</div>
					</div>

					{/* Column 3 */}
					<div className="grid grid-rows-[auto_1fr_auto] gap-8">
						{" "}
						{/* Removed max-w-sm */}
						<h2 className="text-2xl font-semibold">
							Differentiate Early. And Often.
						</h2>
						<p className="text-base font-medium text-gray-600 leading-relaxed">
							Your competitors are constantly evolving, changing prices,
							launching new features. It's not what you know about your
							competitors - it's what you don't know yet that matters.
						</p>
						<div>
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
	);
};

export default LastOne;
