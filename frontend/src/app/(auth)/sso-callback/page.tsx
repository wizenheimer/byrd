"use client";

import LoadingStep from "@/app/(onboarding)/components/steps/LoadingStep";
import { AuthenticateWithRedirectCallback } from "@clerk/nextjs";

export default function SSOCallback() {
	return (
		<>
			<LoadingStep message="Finishing up..." />
			<AuthenticateWithRedirectCallback
				signInFallbackRedirectUrl="/onboarding"
				signUpFallbackRedirectUrl="/onboarding"
			/>
		</>
	);
}
