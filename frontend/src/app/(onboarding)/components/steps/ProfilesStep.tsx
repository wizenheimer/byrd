"use client";

import { INITIAL_PROFILES } from "@/app/constants/onboarding";
import { useOnboardingActions, useProfiles } from "@/app/store/onboarding";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { cn } from "@/lib/utils";

interface ProfilesStepProps {
	onNext: () => void;
}

export default function ProfileStep({ onNext }: ProfilesStepProps) {
	const profiles = useProfiles();
	const { setProfiles } = useOnboardingActions();

	const toggleProfiles = (title: string) => {
		const updatedProfiles = profiles.includes(title)
			? profiles.filter((profileTitle) => profileTitle !== title)
			: [...profiles, title];

		setProfiles(updatedProfiles);
	};

	return (
		<div className="space-y-6">
			{INITIAL_PROFILES.map((profile) => (
				<div key={profile.id} className="flex items-center space-x-4">
					<Switch
						id={profile.id}
						checked={profiles.includes(profile.title)}
						onCheckedChange={() => toggleProfiles(profile.title)}
						className="data-[state=checked]:bg-blue-600"
					/>
					<div className="flex-1 space-y-1">
						<label
							htmlFor={profile.id}
							className="text-base font-semibold leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
						>
							{profile.title}
						</label>
						<p className="text-sm text-muted-foreground">
							{profile.description}
						</p>
					</div>
				</div>
			))}

			<div className="space-y-6">
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
		</div>
	);
}
