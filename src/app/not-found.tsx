import { Button } from "@/components/ui/button";
import Link from "next/link";

const errorMessages = [
	{
		message:
			"Oops! This page pulled a vanishing act, but don't worry - your competitors aren't going anywhere anytime soon. We're still tracking every move they make.",
	},
	{
		message:
			"404: Seems we lost this page in the shuffle, rest assured - we haven't lost sight of your competitors.",
	},
	{
		message:
			"404: Think of this page as our one blind spots. But don't worry - our competitive tracking is still spot-on.",
	},
	{
		message:
			"404: Well this is embarrassing... but hey, don't worry we haven't lost track of your competitors!",
	},
	{
		message:
			"404: We dropped the ball on this page. But don't worry haven't lost track of your competitors!",
	},
	{
		message:
			"404: Page not found - but found something better: your competitors' latest moves! Check them out now.",
	},
	{
		message:
			"404: Seems this page is an uncharted territory. Rest assured, we're still mapping your competitors' every move.",
	},
	{
		message:
			"Oops! Dead end here. But rest assured, competitor tracking is full of live leads!",
	},
	{
		message:
			"404: This page is playing hard to get. Your competitors? We've got them right where we can see them.",
	},
	{
		message:
			"404: Looks like we misplaced this page. Good thing we never misplace your competitor data! Go check it out.",
	},
	{
		message:
			"404: This page is a ghost town. But don't worry, we're still tracking your competitors' every move.",
	},
	{
		message:
			"Oops! Wrong turn here. But don't worry, your market intelligence is still heading in the right direction.",
	},
];

export default function NotFound() {
	const randomMessage =
		errorMessages[Math.floor(Math.random() * errorMessages.length)];

	return (
		<div className="w-full bg-background relative overflow-hidden">
			<div className="max-w-7xl mx-auto px-4 py-24">
				{/* Error Content */}
				<div className="text-center max-w-3xl mx-auto">
					<h1 className="text-5xl font-bold tracking-tight mb-6">Whoops</h1>
					<p className="text-lg text-gray-600 mb-8">{randomMessage.message}</p>
					<div className="flex gap-4 justify-center">
						<Link href="/">
							<Button className="bg-black text-white hover:bg-black/90 px-8">
								Back to Homepage
							</Button>
						</Link>
						<Link href="/help">
							<Button variant="outline" className="border-gray-200">
								Help Center
							</Button>
						</Link>
					</div>
				</div>
			</div>
		</div>
	);
}
