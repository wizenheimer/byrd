"use client";

import LoadingStep from "@/app/(onboarding)/components/steps/LoadingStep";
import {
  type OnboardingData,
  persistOnboardingData,
} from "@/app/actions/onboarding";
import {
  type OnboardingState,
  useOnboardingStore,
} from "@/app/store/onboarding";
import { useAuth, useUser } from "@clerk/nextjs";
import { useRouter } from "next/navigation";
import { useEffect, useRef } from "react";

function transformOnboardingData(stateData: OnboardingState): OnboardingData {
  return {
    competitors: stateData.competitors,
    profiles: stateData.profiles,
    features: stateData.features,
    team: stateData.team,
  };
}

export default function OnboardingComplete() {
  const { isLoaded, isSignedIn, user } = useUser();
  const { getToken } = useAuth();
  const router = useRouter();
  const onboardingState = useOnboardingStore();
  const persistAttemptedRef = useRef(false);

  useEffect(() => {
    // If we've already attempted to persist, don't try again
    if (persistAttemptedRef.current) {
      return;
    }

    const persistData = async () => {
      if (!isLoaded) return;

      if (!isSignedIn) {
        onboardingState.reset();
        router.push("/");
        return;
      }

      if (!user) return;

      try {
        // Mark that we've attempted persistence
        persistAttemptedRef.current = true;

        const stateData: OnboardingState = {
          currentStep: onboardingState.currentStep,
          competitors: onboardingState.competitors,
          profiles: onboardingState.profiles,
          features: onboardingState.features,
          team: onboardingState.team,
        };

        const token = await getToken();
        if (!token) {
          throw new Error("Failed to retrieve authentication token");
        }

        const onboardingData = transformOnboardingData(stateData);
        const result = await persistOnboardingData(onboardingData, token);

        if (result.success) {
          onboardingState.reset();
          router.push("/waitlist");
        } else {
          throw new Error("Failed to persist onboarding data");
        }
      } catch (error) {
        console.error("Error persisting onboarding data:", error);
        // Reset the attempt flag on error so user can try again
        persistAttemptedRef.current = false;
      }
    };

    persistData();
  }, [isLoaded, isSignedIn, user, getToken, onboardingState, router]);

  return (
    <LoadingStep
      message={persistAttemptedRef.current ? "Saving your preferences..." : "Almost there..."}
    />
  );
}
