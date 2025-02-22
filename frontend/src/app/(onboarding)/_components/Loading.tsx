// app/(onboarding)/_components/Loading.tsx
import { Loader2 } from "lucide-react";

export default function Loading() {
	return (
		<div className="flex min-h-screen items-center justify-center">
			<div className="text-center space-y-4">
				<Loader2 className="h-8 w-8 animate-spin mx-auto" />
				<h2 className="text-xl font-semibold">
					Completing Slack Installation...
				</h2>
				<p className="text-gray-600">
					Please wait while we finish setting up your Slack integration.
				</p>
			</div>
		</div>
	);
}
