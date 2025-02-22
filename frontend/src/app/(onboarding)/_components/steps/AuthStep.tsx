// app/(onboarding)/_components/steps/AuthStep.tsx
"use client";

import { useOnboardingStore } from "@/app/store/onboarding";
import { Button } from "@/components/ui/button";
import { ToastAction } from "@/components/ui/toast";
import { useToast } from "@/hooks/use-toast";
import { Slack } from "lucide-react";
import { handleSlackInit } from "../../_actions/onboarding";

export default function AuthStep() {
	const { toast } = useToast();
	const reset = useOnboardingStore((state) => state.reset);
	const competitors = useOnboardingStore.use.competitors();
	const profiles = useOnboardingStore.use.profiles();
	const features = useOnboardingStore.use.features();

	const handleSlackInstall = async () => {
		try {
			const result = await handleSlackInit({
				competitors: competitors,
				features: features,
				profiles: profiles,
			});
			reset();

			if (!result.success) {
				throw new Error(result.error);
			}

			window.location.href = result.oauth_url;
		} catch (error) {
			console.error(error);
			toast({
				variant: "destructive",
				title: "Installation Failed",
				description: "Failed to initiate Slack installation. Please try again.",
				action: <ToastAction altText="Try again">Try again</ToastAction>,
			});
		}
	};

	return (
		<div className="space-y-4">
			<Button
				variant="outline"
				className="relative h-12 w-full justify-center text-base font-normal"
				onClick={handleSlackInstall}
			>
				<div className="absolute left-4 size-5">
					<Slack />
				</div>
				Sign in with Slack
			</Button>
		</div>
	);
}
