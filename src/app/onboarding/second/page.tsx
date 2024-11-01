"use client";

import OnboardingPreviewPane from "@/components/block/OnboardingPreviewPane";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { cn } from "@/lib/utils";
import * as React from "react";

interface FeatureItem {
	id: string
	title: string
	description: string
	enabled: boolean
}

export default function Component() {
	const [features, setFeatures] = React.useState<FeatureItem[]>([
		{ id: "1", title: "Product", description: "Catch product evolution in real-time", enabled: true },
		{ id: "2", title: "Pricing", description: "Never be the last to know about a price war", enabled: false },
		{ id: "3", title: "Partnership", description: "Track who's teaming up with whom", enabled: false },
		{ id: "4", title: "Branding", description: "Monitor messaging shifts, and identity changes", enabled: false },
		{ id: "5", title: "Positioning", description: "Track narratives before they go mainstream", enabled: false },
	])

	const toggleFeature = (id: string) => {
		setFeatures(features.map(feature =>
			feature.id === id ? { ...feature, enabled: !feature.enabled } : feature
		))
	}

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
					<div className="space-y-6">
						{features.map((feature) => (
							<div key={feature.id} className="flex items-center space-x-4">
								<Switch
									id={feature.id}
									checked={feature.enabled}
									onCheckedChange={() => toggleFeature(feature.id)}
								/>
								<div className="flex-1 space-y-1">
									<label
										htmlFor={feature.id}
										className="text-base font-semibold leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
									>
										{feature.title}
									</label>
									<p className="text-sm text-muted-foreground">
										{feature.description}
									</p>
								</div>
							</div>
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
