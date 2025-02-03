import { STORAGE_KEYS } from "@/constants/storage";
import { createSelectors } from "@/lib/utils";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import { devtools } from "zustand/middleware";
import { INITIAL_PROFILES } from "../constants/onboarding";
import { ProfileType } from "../types/competitor_page";

export interface OnboardingState {
  currentStep: number;
  competitors: string[]; // competitor urls
  profiles: ProfileType[]; // profile names
  features: string[]; // feature names
  team: string[]; // team emails
}

export interface OnboardingActions {
  setCurrentStep: (step: number) => void;
  setCompetitors: (competitors: string[]) => void;
  setProfiles: (profiles: ProfileType[]) => void;
  setFeatures: (features: string[]) => void;
  setTeam: (team: string[]) => void;
  reset: () => void;
}

const initialState: OnboardingState = {
  currentStep: 1,
  competitors: [],
  profiles: INITIAL_PROFILES.filter((profile) => profile.enabled).map(
    (profile) => profile.profile_key
  ),
  features: ["Product"],
  team: [],
};

export type OnboardingStore = OnboardingState & OnboardingActions;

const useOnboardingStoreBase = create<OnboardingStore>()(
  devtools(
    persist(
      (set) => ({
        ...initialState,

        setCurrentStep: (step) =>
          set((state) => ({ ...state, currentStep: step })),

        setCompetitors: (competitors) =>
          set((state) => ({ ...state, competitors: competitors })),

        setProfiles: (profiles) =>
          set((state) => ({ ...state, profiles: profiles })),

        setFeatures: (features) =>
          set((state) => ({ ...state, features: features })),

        setTeam: (team) => set((state) => ({ ...state, team: team })),

        reset: () => set(initialState),
      }),
      {
        name: STORAGE_KEYS.FORM_DATA, // Storage key
        storage: createJSONStorage(() => localStorage),
        partialize: (state) => ({
          // Only persist these state properties
          currentStep: state.currentStep,
          competitors: state.competitors,
          profiles: state.profiles,
          features: state.features,
          team: state.team,
        }),
      }
    ),
    {
      name: "onboarding-store",
    }
  )
);

export const useOnboardingStore = createSelectors(useOnboardingStoreBase);

// Note: deprecated hooks with utility functions
// Selector hooks for better performance
// export const useCurrentStep = () =>
// useOnboardingStore((state) => state.currentStep);
// export const useCompetitors = () =>
// useOnboardingStore((state) => state.competitors);
// export const useProfiles = () => useOnboardingStore((state) => state.profiles);
// export const useFeatures = () => useOnboardingStore((state) => state.features);
// export const useTeam = () => useOnboardingStore((state) => state.team);

// Action selector hooks
// export const useOnboardingActions = () => {
//   const store = useOnboardingStore();
//   return {
//     setCurrentStep: store.setCurrentStep,
//     setCompetitors: store.setCompetitors,
//     setProfiles: store.setProfiles,
//     setFeatures: store.setFeatures,
//     setTeam: store.setTeam,
//     reset: store.reset,
//   };
// };
