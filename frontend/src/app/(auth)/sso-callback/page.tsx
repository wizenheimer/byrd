// src/app/(onboarding)/sso-callback/page.tsx
"use client";

import { AuthenticateWithRedirectCallback } from "@clerk/nextjs";
import LoadingStep from "@/app/(auth)/components/steps/LoadingStep";

export default function SSOCallback() {
  return (
    <>
      <LoadingStep message="Finishing up..." />
      <AuthenticateWithRedirectCallback
        signInFallbackRedirectUrl="/onboarding-complete"
        signUpFallbackRedirectUrl="/onboarding-complete"
      />
    </>
  );
}
