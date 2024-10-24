import Link from "next/link";
import React from "react";
import PreviewContainer from "./PreviewContainer";

const LastOne = () => {
	return (
		<div className="w-full bg-background relative overflow-hidden mt-40">
			{/* Hero Text Content */}
			<div className="max-w-7xl mx-auto px-4 pt-20 pb-20">
				<div className="text-center max-w-4xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-4">
						Don&apos;t be the last one to know.
					</h1>
					<p className="text-lg text-gray-600 mt-8">
						Spot emerging competitors long before they become existential
						threats.
						<br />
						No more of &apos;How did we miss that?&apos; moments.
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
								imageSrc="/assets/shared-view-left.png"
								variant="compact"
								position="right"
							/>
						</div>
						<div className="space-y-6">
							<h2 className="text-3xl font-bold">Shared competitive picture</h2>
							<p className="text-gray-600">
								Right now, you&apos;ve got Marketing hoarding one set of
								competitor data, Sales clutching another, and Product off to
								something else. By the time this fragmented intel makes its way
								up the chain, competitors have already made their move, and
								you&apos;re playing on the back foot.
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
								imageSrc="/assets/shared-view-right.png"
								variant="compact"
								position="right"
							/>
						</div>
						<div className="space-y-6">
							<h2 className="text-3xl font-bold">Enrich your knowledge base</h2>
							<p className="text-gray-600">
								Existing tools don&apos;t solve data silos; they create new ones
								and leave the plumbing for you. Byrd is built differently.
							</p>
							<p className="text-gray-600">
								Byrd offers first class integrations with the tools you already
								use - Slack, Notion, and even email. This means the competitive
								intelligence you gather doesn&apos;t just sit in yet another
								dashboard you&apos;ll forget to check. Instead, it flows
								seamlessly into your existing knowledge base and communication
								channels.
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
		</div>
	);
};

export default LastOne;
