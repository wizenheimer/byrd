"use client";

import {
	type OnboardingData,
	persistOnboardingData,
} from "@/app/_actions/onboarding";
import { onboardingFormSchema } from "@/app/_types/onboarding";
import LoadingStep from "@/components/steps/LoadingStep";
import { STORAGE_KEYS } from "@/constants/storage";
import { useUser } from "@clerk/nextjs";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

function loadOnboardingDataFromStorage(
	storageKey: string,
	// biome-ignore lint/suspicious/noExplicitAny: <explanation>
	userData: any,
): OnboardingData {
	// Get data from localStorage
	const rawData = localStorage.getItem(storageKey);
	const parsedData = rawData ? JSON.parse(rawData) : {};

	// Parse and validate the form data
	const formData = onboardingFormSchema.parse(parsedData);

	// Transform to OnboardingData format
	return {
		clerkId: userData.id,
		email: userData.primaryEmailAddress?.emailAddress || "",
		firstName: userData.firstName || "",
		lastName: userData.lastName || "",
		competitors: formData.competitors.map((competitor) => ({
			url: competitor.url,
		})),
		features: formData.features
			.filter((feature) => feature.enabled)
			.map((feature) => ({
				title: feature.title,
			})),
		channels: formData.channels.map((channel) => ({
			title: channel,
		})),
		team: formData.team.flatMap((teamForm) =>
			teamForm.members.map((member) => ({
				email: member.email,
			})),
		),
	};
}

export default function OnboardingComplete() {
	const { isLoaded, isSignedIn, user } = useUser();
	const router = useRouter();
	const [isSubmitting, setIsSubmitting] = useState(false);

	useEffect(() => {
		const persistData = async () => {
			if (isLoaded && isSignedIn && user) {
				setIsSubmitting(true);
				try {
					const onboardingData = loadOnboardingDataFromStorage(
						STORAGE_KEYS.FORM_DATA,
						user,
					);

					await persistOnboardingData(onboardingData);

					// Clear local storage
					localStorage.removeItem(STORAGE_KEYS.STEP);
					localStorage.removeItem(STORAGE_KEYS.FORM_DATA);

					router.push("/waitlist");
				} catch (error) {
					console.error("Error persisting onboarding data:", error);
					// Handle error (e.g., show error message to user)
				} finally {
					setIsSubmitting(false);
				}
			} else if (isLoaded && !isSignedIn) {
				// Clear local storage
				localStorage.removeItem(STORAGE_KEYS.STEP);
				localStorage.removeItem(STORAGE_KEYS.FORM_DATA);

				// Redirect to home page
				router.push("/");
			}
		};

		persistData();
	}, [isLoaded, isSignedIn, user, router]);

	if (!isLoaded || !isSignedIn || isSubmitting) {
		return <LoadingStep />;
	}

	return <LoadingStep />;
}
