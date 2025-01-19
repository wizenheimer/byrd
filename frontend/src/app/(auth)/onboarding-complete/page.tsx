// src/app/(onboarding)/onboarding-complete/page.tsx
"use client";

import { OnboardingData, persistOnboardingData } from "@/app/_actions/onboarding";
import LoadingStep from "@/app/(auth)/components/steps/LoadingStep";
import { OnboardingState, useOnboardingStore } from "@/app/_store/onboarding";
import { useUser } from "@clerk/nextjs";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

function transformOnboardingData(
  stateData: OnboardingState
): OnboardingData {
  return {
    competitors: stateData.competitors,
    features: stateData.enabledFeatures,
    channels: stateData.channels,
    team: stateData.team,
  };
}

export default function OnboardingComplete() {
  const { isLoaded, isSignedIn, user } = useUser();
  const router = useRouter();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const onboardingState = useOnboardingStore();

  useEffect(() => {
    const persistData = async () => {
      if (isLoaded && isSignedIn && user) {
        setIsSubmitting(true);
        try {
          // Pick only the state properties, excluding actions
          const stateData: OnboardingState = {
            currentStep: onboardingState.currentStep,
            competitors: onboardingState.competitors,
            enabledFeatures: onboardingState.enabledFeatures,
            channels: onboardingState.channels,
            team: onboardingState.team,
          };

          const onboardingData = transformOnboardingData(stateData);
          const result = await persistOnboardingData(onboardingData);

          if (result.success) {
            // Clear onboarding state
            onboardingState.reset();
            router.push("/waitlist");
          } else {
            throw new Error("Failed to persist onboarding data");
          }
        } catch (error) {
          console.error("Error persisting onboarding data:", error);
          // You might want to show an error message to the user here
        } finally {
          setIsSubmitting(false);
        }
      } else if (isLoaded && !isSignedIn) {
        // Clear state and redirect if not signed in
        onboardingState.reset();
        router.push("/");
      }
    };

    persistData();
  }, [isLoaded, isSignedIn, user, router, onboardingState]);

  return (
    <LoadingStep message={isSubmitting ? "Saving your preferences..." : "Almost there..."} />
  );
}
