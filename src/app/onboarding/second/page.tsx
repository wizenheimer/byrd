"use client";

import { OnboardingHeader } from "@/components/OnboardingHeader";
import { OnboardingLayout } from "@/components/OnboardingLayout";
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

export default function FeaturesStep() {
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
		<OnboardingLayout
			previewImage="/onboarding/second.png"
			previewAlt="Dashboard Preview"
		>
			<OnboardingHeader
				title="Measure What Matters"
				description="Choose your signals. Cut through the noise."
			/>

			{/* Property Selector */}
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

			{/* Property Footnote */}
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
		</OnboardingLayout>
	);
}
