import { STORAGE_KEYS } from "@/constants/storage";
import { createSelectors } from "@/lib/utils";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import { devtools } from "zustand/middleware";
import { INITIAL_PROFILES } from "../(onboarding)/_constants/onboarding";
import { ProfileType } from "../(onboarding)/_types/onboarding";

export interface OnboardingState {
	currentStep: number;
	competitors: string[]; // competitor urls
	profiles: ProfileType[]; // profile names
	features: string[]; // feature names
}

export interface OnboardingActions {
	setCurrentStep: (step: number) => void;
	setCompetitors: (competitors: string[]) => void;
	setProfiles: (profiles: ProfileType[]) => void;
	setFeatures: (features: string[]) => void;
	reset: () => void;
}

const initialState: OnboardingState = {
	currentStep: 1,
	competitors: [],
	profiles: INITIAL_PROFILES.filter((profile) => profile.enabled).map(
		(profile) => profile.profile_key,
	),
	features: ["Product"],
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
				}),
			},
		),
		{
			name: "onboarding-store",
		},
	),
);

export const useOnboardingStore = createSelectors(useOnboardingStoreBase);
