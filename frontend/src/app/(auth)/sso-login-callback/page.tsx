"use client";

import LoadingStep from "@/app/(onboarding)/components/steps/LoadingStep";
import { AuthenticateWithRedirectCallback } from "@clerk/nextjs";

export default function SSOLoginCallback() {
	return (
		<>
			<LoadingStep message="Completing authentication..." />
			<AuthenticateWithRedirectCallback
				signInFallbackRedirectUrl="/dashboard"
				signUpFallbackRedirectUrl="/dashboard"
			/>
		</>
	);
}
