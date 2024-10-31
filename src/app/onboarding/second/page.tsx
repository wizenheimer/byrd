"use client";

import OnboardingPreviewPane from "@/components/block/OnboardingPreviewPane";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import {
	// ArrowLeft,
	// BarChart2,
	Box,
	DollarSign,
	Lightbulb,
	Link2,
	Palette,
	Share2,
} from "lucide-react";
// import Link from "next/link";
import * as React from "react";

export default function Component() {
	const [selected, setSelected] = React.useState<string[]>([
		"pricing",
		"product",
		"positioning",
	]);

	const cards = [
		{ id: "pricing", icon: DollarSign, label: "Pricing" },
		{ id: "product", icon: Box, label: "Product" },
		{ id: "branding", icon: Palette, label: "Branding" },
		{ id: "positioning", icon: Lightbulb, label: "Positioning" },
		{ id: "integrations", icon: Link2, label: "Integrations" },
		{ id: "partnership", icon: Share2, label: "Partnership" },
	];

	const toggleCard = (id: string) => {
		setSelected((prev) =>
			prev.includes(id) ? prev.filter((item) => item !== id) : [...prev, id],
		);
	};

	return (
		<div className="flex min-h-screen flex-col lg:flex-row">
			{/* Left Side - Signal Selection */}
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
					<div className="grid grid-cols-2 gap-4">
						{cards.map(({ id, icon: Icon, label }) => (
							<button
								key={id}
								onClick={() => toggleCard(id)}
								type="button"
								className={`relative flex h-32 flex-col items-center justify-center rounded-xl border-2 transition-colors
				  ${
						selected.includes(id)
							? "border-primary bg-primary/5"
							: "border-border bg-background hover:border-primary/50"
					}`}
							>
								{selected.includes(id) && (
									<div className="absolute right-2 top-2 size-4 rounded-full bg-primary" />
								)}
								<Icon className="mb-2 size-6" />
								<span className="text-sm font-medium">{label}</span>
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
				imageSrc="/onboarding/second.png"
				altText="Dashboard Preview"
			/>
		</div>
	);
}
