// src/app/(onboarding)/sso-callback/page.tsx
"use client";

import { AuthenticateWithRedirectCallback } from "@clerk/nextjs";

export default function SSOCallback() {
  return (
    <AuthenticateWithRedirectCallback
      signInFallbackRedirectUrl="/onboarding-complete"
      signUpFallbackRedirectUrl="/onboarding-complete"
    />
  );
}
