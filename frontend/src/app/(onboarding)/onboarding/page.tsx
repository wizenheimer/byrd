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
import { ToastAction } from "@/components/ui/toast";
import { useToast } from "@/hooks/use-toast";
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
  const retryCountRef = useRef(0);
  const { toast } = useToast();

  useEffect(() => {
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

      persistAttemptedRef.current = true;

      const attemptPersist = async () => {
        try {
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
            return true;
          }
          throw new Error("Failed to persist onboarding data");
        } catch (error) {
          console.error(`Persistence attempt ${retryCountRef.current + 1} failed:`, error);
          return false;
        }
      };

      while (retryCountRef.current < 3) {
        const success = await attemptPersist();
        if (success) return;

        retryCountRef.current += 1;
        if (retryCountRef.current < 3) {
          await new Promise(resolve => setTimeout(resolve, 1000));
        }
      }

      // All retries failed
      persistAttemptedRef.current = false;
      toast({
        title: "Uh oh! Something went wrong.",
        description: "We couldn't get you onboarded.",
        action: (
          <ToastAction altText="Try Again" onClick={() => {
            onboardingState.reset();
            router.push("/");
          }}>
            Go to Homepage
          </ToastAction>
        ),
      });
    };

    persistData();
  }, [isLoaded, isSignedIn, user, getToken, onboardingState, router, toast]);

  return (
    <LoadingStep
      message={persistAttemptedRef.current ? "Saving your preferences..." : "Almost there..."}
    />
  );
}
