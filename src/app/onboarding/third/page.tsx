"use client";
import OnboardingPreviewPane from "@/components/block/OnboardingPreviewPane";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import {
	// ArrowLeft,
	// BarChart2,
	Inbox,
	Megaphone,
	Rss,
	Share2,
} from "lucide-react";
// import Link from "next/link";
import * as React from "react";

export default function Component() {
	const [selected, setSelected] = React.useState<string[]>([
		"inbox",
		"mentions",
	]);

	const cards = [
		{
			id: "inbox",
			icon: Inbox,
			title: "Inbox",
			description: "Monitor the direct line to their customers",
		},
		{
			id: "social",
			icon: Share2,
			title: "Social",
			description: "Follow their social playbook as it unfolds",
		},
		{
			id: "mentions",
			icon: Megaphone,
			title: "Mentions",
			description: "Beat them to their own announcement",
		},
		{
			id: "content",
			icon: Rss,
			title: "Content",
			description: "Watch their messaging evolve, post by post",
		},
	];

	const toggleCard = (id: string) => {
		setSelected((prev) =>
			prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id],
		);
	};

	return (
		<div className="flex min-h-screen flex-col lg:flex-row">
			{/* Left Side - Channel Selection */}
			<div className="flex flex-1 flex-col bg-white p-8 lg:p-12">
				<div className="mb-16">
					<span className="text-xl font-semibold">byrd</span>
				</div>
				<div className="mx-auto w-full max-w-[440px] space-y-12">
					<div className="space-y-3">
						<h1 className="text-3xl font-bold tracking-tight">
							Measure What Matters
						</h1>
						<p className="text-base text-muted-foreground">
							Choose your signals. Cut through the noise.
						</p>
					</div>
					<div className="space-y-6">
						{cards.map(({ id, icon: Icon, title, description }) => (
							<button
								key={id}
								type="button"
								onClick={() => toggleCard(id)}
								className={`relative flex w-full items-start gap-4 rounded-xl border-2 p-4 text-left transition-colors
				  ${
						selected.includes(id)
							? "border-primary bg-primary/5"
							: "border-border bg-background hover:border-primary/50"
					}`}
							>
								<Icon className="mt-1 size-5 shrink-0" />
								<div className="space-y-1">
									<div className="font-medium">{title}</div>
									<div className="text-sm text-muted-foreground">
										{description}
									</div>
								</div>
								{selected.includes(id) && (
									<div className="absolute right-3 top-1/2 size-4 -translate-y-1/2 rounded-full bg-primary" />
								)}
							</button>
						))}
					</div>
					<div className="space-y-6">
						<Button
							className={cn(
								"w-full h-12 text-base font-semibold",
								"bg-[#14171F] hover:bg-[#14171F]/90",
								"transition-colors duration-200",
								"disabled:opacity-50 disabled:cursor-not-allowed",
							)}
							size="lg"
						>
							Continue
						</Button>
						<p className="text-sm text-muted-foreground text-center">
							You can always customize them later
						</p>
					</div>
				</div>
			</div>

			{/* Right Side - Dashboard Preview */}
			<OnboardingPreviewPane
				imageSrc="/onboarding/third.png"
				altText="Dashboard Preview"
			/>
		</div>
	);
}
