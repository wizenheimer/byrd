// app/(onboarding)/_components/steps/AuthStep.tsx
import { Button } from "@/components/ui/button";
import { ToastAction } from "@/components/ui/toast";
import { useToast } from "@/hooks/use-toast";
import { Slack } from "lucide-react";

export default function AuthStep() {
	const { toast } = useToast();

	const handleSlackInstall = async () => {
		try {
			const response = await fetch("/api/oauth/slack/init", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify({
					competitors: [],
					features: [],
					profiles: [],
				}),
			});

			if (!response.ok) {
				throw new Error("Failed to get OAuth URL");
			}

			const { oauth_url } = await response.json();
			window.location.href = oauth_url;
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
