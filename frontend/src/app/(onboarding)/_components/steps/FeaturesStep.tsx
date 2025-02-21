"use client";

import type { FeatureCard } from "@/app/(onboarding)/_schema/onboarding";
import { useOnboardingStore } from "@/app/store/onboarding";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface FeaturesStepProps {
	cards: FeatureCard[];
	onNext: () => void;
}

export default function FeaturesStep({ cards, onNext }: FeaturesStepProps) {
	const selectedFeatures = useOnboardingStore.use.features();
	const setFeatures = useOnboardingStore.use.setFeatures();

	const toggleFeature = (title: string) => {
		const newFeatures = selectedFeatures.includes(title)
			? selectedFeatures.filter((featureTitle) => featureTitle !== title)
			: [...selectedFeatures, title];

		setFeatures(newFeatures);
	};

	return (
		<div className="space-y-6">
			{cards.map(({ id, icon, title, description }) => (
				<FeatureButton
					key={id}
					id={id}
					icon={icon}
					title={title}
					description={description}
					isSelected={selectedFeatures.includes(title)}
					onClick={() => toggleFeature(title)}
				/>
			))}
			<Button
				className={cn(
					"w-full h-12 text-base font-semibold",
					"bg-[#14171F] hover:bg-[#14171F]/90",
					"transition-colors duration-200",
					"disabled:opacity-50 disabled:cursor-not-allowed",
				)}
				size="lg"
				onClick={onNext}
			>
				Continue
			</Button>
			<p className="text-sm text-muted-foreground text-center">
				You can always customize them later
			</p>
		</div>
	);
}

const FeatureButton = ({
	// id,
	icon: Icon,
	title,
	description,
	isSelected,
	onClick,
}: {
	id: string;
	icon: React.ComponentType<{ className?: string }>;
	title: string;
	description: string;
	isSelected: boolean;
	onClick: () => void;
}) => (
	<button
		type="button"
		onClick={onClick}
		className={cn(
			"relative flex w-full items-start gap-4 rounded-xl border-2 p-4 text-left transition-colors",
			isSelected
				? "border-primary bg-primary/5"
				: "border-border bg-background hover:border-primary/50",
		)}
	>
		<Icon className="mt-1 size-5 shrink-0" />
		<div className="space-y-1">
			<div className="font-medium">{title}</div>
			<div className="text-sm text-muted-foreground">{description}</div>
		</div>
		{isSelected && (
			<div className="absolute right-3 top-1/2 size-4 -translate-y-1/2 rounded-full bg-primary" />
		)}
	</button>
);
