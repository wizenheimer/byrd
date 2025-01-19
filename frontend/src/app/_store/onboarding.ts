// src/stores/onboarding.ts
import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { devtools } from "zustand/middleware";
import { STORAGE_KEYS } from "@/constants/storage";
import { INITIAL_FEATURES } from "../_constants/onboarding";

export interface OnboardingState {
  currentStep: number;
  competitors: string[]; // Just URLs
  enabledFeatures: string[]; // feature names
  channels: string[]; // channel names
  team: string[]; // Just emails
}

export interface OnboardingActions {
  setCurrentStep: (step: number) => void;
  setCompetitors: (competitors: string[]) => void;
  setEnabledFeatures: (features: string[]) => void;
  setChannels: (channels: string[]) => void;
  setTeam: (team: string[]) => void;
  reset: () => void;
}

const initialState: OnboardingState = {
  currentStep: 1,
  competitors: [],
  enabledFeatures: INITIAL_FEATURES.filter((feature) => feature.enabled).map(
    (feature) => feature.id
  ),
  channels: ["product"],
  team: [],
};

export type OnboardingStore = OnboardingState & OnboardingActions;

export const useOnboardingStore = create<OnboardingStore>()(
  devtools(
    persist(
      (set) => ({
        ...initialState,

        setCurrentStep: (step) =>
          set((state) => ({ ...state, currentStep: step })),

        setCompetitors: (competitors) =>
          set((state) => ({ ...state, competitors })),

        setEnabledFeatures: (features) =>
          set((state) => ({ ...state, enabledFeatures: features })),

        setChannels: (channels) => set((state) => ({ ...state, channels })),

        setTeam: (team) => set((state) => ({ ...state, team })),

        reset: () => set(initialState),
      }),
      {
        name: STORAGE_KEYS.FORM_DATA, // Storage key
        storage: createJSONStorage(() => localStorage),
        partialize: (state) => ({
          // Only persist these state properties
          currentStep: state.currentStep,
          competitors: state.competitors,
          enabledFeatures: state.enabledFeatures,
          channels: state.channels,
          team: state.team,
        }),
      }
    ),
    {
      name: "onboarding-store",
    }
  )
);

// Selector hooks for better performance
export const useCurrentStep = () =>
  useOnboardingStore((state) => state.currentStep);
export const useCompetitors = () =>
  useOnboardingStore((state) => state.competitors);
export const useEnabledFeatures = () =>
  useOnboardingStore((state) => state.enabledFeatures);
export const useChannels = () => useOnboardingStore((state) => state.channels);
export const useTeam = () => useOnboardingStore((state) => state.team);

// Action selector hooks
export const useOnboardingActions = () => {
  const store = useOnboardingStore();
  return {
    setCurrentStep: store.setCurrentStep,
    setCompetitors: store.setCompetitors,
    setEnabledFeatures: store.setEnabledFeatures,
    setChannels: store.setChannels,
    setTeam: store.setTeam,
    reset: store.reset,
  };
};
