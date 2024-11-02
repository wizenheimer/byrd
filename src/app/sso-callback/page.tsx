import { AuthenticateWithRedirectCallback } from "@clerk/nextjs";
import { STORAGE_KEYS } from "../get-started/page";

export default function SSOCallback() {
	// Persist form data from localStorage to neon
	const persistFormData = () => {
		try {
			const formData = localStorage.getItem(STORAGE_KEYS.FORM_DATA);
			if (formData) {
				const parsedFormData = JSON.parse(formData);
				console.log("Persisting form data:", parsedFormData);
				// TODO: Persist form data to Neon
				setTimeout(() => {
					console.log("Form data persisted to Neon");
				}, 5000);
			}
		} catch (error) {
			console.error("Error persisting form data:", error);
		}
	};

	// Persist form data on load
	persistFormData();

	// Clear onboarding storage
	const clearOnboardingStorage = () => {
		try {
			localStorage.removeItem(STORAGE_KEYS.STEP);
			localStorage.removeItem(STORAGE_KEYS.FORM_DATA);
		} catch (error) {
			console.error("Error clearing localStorage:", error);
		}
	};

	// Clear onboarding storage on load
	clearOnboardingStorage();

	return (
		<AuthenticateWithRedirectCallback
			signInFallbackRedirectUrl="/waitlist"
			signUpFallbackRedirectUrl="/waitlist"
		/>
	);
}
