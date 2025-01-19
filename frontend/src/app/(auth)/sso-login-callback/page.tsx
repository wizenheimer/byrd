// src/app/(onboarding)/sso-login-callback/page.tsx
"use client";

import { AuthenticateWithRedirectCallback } from "@clerk/nextjs";
import LoadingStep from "@/app/(auth)/components/steps/LoadingStep";

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
