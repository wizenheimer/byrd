"use client";

import { Button } from "@/components/ui/button";
import { AnimatePresence, motion } from "framer-motion";
import { ChevronLeft, Inbox, Megaphone, Rss, Share2 } from "lucide-react";
import { useEffect, useState } from "react";
import AuthStep from "../../components/steps/AuthStep";
import ChannelsStep from "../../components/steps/ChannelsStep";
import CompetitorStep from "../../components/steps/CompetitorStep";
import FeaturesStep from "../../components/steps/FeaturesStep";
import TeamStep from "../../components/steps/TeamStep";
import type { ChannelCard, OnboardingFormData } from "../_types/onboarding";

const STORAGE_KEYS = {
	STEP: "onboarding_current_step",
	FORM_DATA: "onboarding_form_data",
} as const;

interface Step {
	title: string;
	description: string;
	image: string;
}

const initialFormStep = 1;
const initialFormState: OnboardingFormData = {
	competitors: [],
	features: [
		{
			id: "1",
			title: "Product",
			description: "Catch product evolution in real-time",
			enabled: true,
		},
		{
			id: "2",
			title: "Pricing",
			description: "Never be the last to know about a price war",
			enabled: false,
		},
		{
			id: "3",
			title: "Partnership",
			description: "Track who's teaming up with whom",
			enabled: false,
		},
		{
			id: "4",
			title: "Branding",
			description: "Monitor messaging shifts, and identity changes",
			enabled: false,
		},
		{
			id: "5",
			title: "Positioning",
			description: "Track narratives before they go mainstream",
			enabled: false,
		},
	],
	channels: ["inbox", "mentions"],
	team: [],
};

const steps: Record<number, Step> = {
	1: {
		title: "Your Market, Your Rules",
		description: "Pick your targets. Add up to 5 competitors.",
		image: "/onboarding/first.png",
	},
	2: {
		title: "Measure What Matters",
		description: "Choose your signals. Cut through the noise.",
		image: "/onboarding/second.png",
	},
	3: {
		title: "Never Miss A Beat",
		description: "Your competitors are everywhere. So are we.",
		image: "/onboarding/third.png",
	},
	4: {
		title: "Build Your War Room",
		description: "Business is a team sport. Bring in your heavy hitters.",
		image: "/onboarding/four.png",
	},
	5: {
		title: "You're almost there",
		description: "Quick auth, then let's get started.",
		image: "/onboarding/five.png",
	},
};

const channelCards: ChannelCard[] = [
	{
		id: "inbox",
		icon: Inbox,
		title: "Inbox",
		description: "Monitor the direct line to their customers",
	},
	{
		id: "social",
		icon: Share2,
		title: "Social",
		description: "Follow their social playbook as it unfolds",
	},
	{
		id: "mentions",
		icon: Megaphone,
		title: "Mentions",
		description: "Beat them to their own announcement",
	},
	{
		id: "content",
		icon: Rss,
		title: "Content",
		description: "Watch their messaging evolve, post by post",
	},
];

const MultiStepOnboarding = () => {
	const [currentStep, setCurrentStep] = useState(initialFormStep);
	const [formData, setFormData] =
		useState<OnboardingFormData>(initialFormState);
	const [isClient, setIsClient] = useState(false);

	useEffect(() => {
		setIsClient(true);
		try {
			const savedStep = localStorage.getItem(STORAGE_KEYS.STEP);
			const savedData = localStorage.getItem(STORAGE_KEYS.FORM_DATA);

			if (savedStep) {
				setCurrentStep(Number.parseInt(savedStep, 10));
			}

			if (savedData) {
				setFormData(JSON.parse(savedData));
			}
		} catch (error) {
			console.error("Error reading from localStorage:", error);
		}
	}, []);

	useEffect(() => {
		if (isClient) {
			try {
				localStorage.setItem(STORAGE_KEYS.STEP, currentStep.toString());
				localStorage.setItem(STORAGE_KEYS.FORM_DATA, JSON.stringify(formData));
			} catch (error) {
				console.error("Error saving to localStorage:", error);
			}
		}
	}, [currentStep, formData, isClient]);

	const clearOnboardingStorage = () => {
		try {
			localStorage.removeItem(STORAGE_KEYS.STEP);
			localStorage.removeItem(STORAGE_KEYS.FORM_DATA);
		} catch (error) {
			console.error("Error clearing localStorage:", error);
		}
	};

	const handleNext = () => {
		if (currentStep < 5) {
			setCurrentStep((prev) => prev + 1);
		}
	};

	const handleBack = () => {
		if (currentStep > 1) {
			setCurrentStep((prev) => prev - 1);
		}
	};

	const handleComplete = () => {
		console.log("cleaning up onboarding state");
		clearOnboardingStorage();
		// Add your completion logic here
	};

	const renderStep = () => {
		switch (currentStep) {
			case 1:
				return (
					<CompetitorStep
						formData={formData}
						setFormData={setFormData}
						onNext={handleNext}
					/>
				);
			case 2:
				return (
					<FeaturesStep
						formData={formData}
						setFormData={setFormData}
						onNext={handleNext}
					/>
				);
			case 3:
				return (
					<ChannelsStep
						formData={formData}
						setFormData={setFormData}
						onNext={handleNext}
						cards={channelCards}
					/>
				);
			case 4:
				return (
					<TeamStep
						formData={formData}
						setFormData={setFormData}
						onNext={handleNext}
					/>
				);
			case 5:
				return <AuthStep onComplete={handleComplete} />;
			default:
				return null;
		}
	};

	if (!isClient) {
		return null; // or a loading spinner
	}

	return (
		<div className="flex min-h-screen flex-col lg:flex-row">
			<div className="flex flex-1 flex-col bg-white p-8 lg:p-12">
				<nav className="mb-16 flex items-center justify-between">
					<div className="flex items-center">
						<AnimatePresence mode="wait">
							{currentStep > 1 && (
								<motion.div
									initial={{ opacity: 0, width: 0 }}
									animate={{ opacity: 1, width: "auto" }}
									exit={{ opacity: 0, width: 0 }}
									transition={{ duration: 0.2 }}
								>
									<Button
										variant="ghost"
										onClick={handleBack}
										className="group -ml-3 h-9 px-2 text-muted-foreground hover:text-foreground"
									>
										<ChevronLeft className="h-4 w-4 transition-transform group-hover:-translate-x-1" />
									</Button>
								</motion.div>
							)}
						</AnimatePresence>
						<span className="text-xl font-semibold">byrd</span>
					</div>
					<div className="text-sm text-muted-foreground">
						Step {currentStep} of {5}
					</div>
				</nav>

				<div className="mx-auto w-full max-w-[440px] space-y-12">
					<AnimatePresence mode="wait">
						<motion.div
							key={currentStep}
							initial={{ opacity: 0, y: 20 }}
							animate={{ opacity: 1, y: 0 }}
							exit={{ opacity: 0, y: -20 }}
							transition={{ duration: 0.3 }}
							className="space-y-6"
						>
							<div className="space-y-3">
								<h1 className="text-2xl font-bold tracking-tight">
									{steps[currentStep].title}
								</h1>
								<p className="text-base text-muted-foreground">
									{steps[currentStep].description}
								</p>
							</div>
							{renderStep()}
						</motion.div>
					</AnimatePresence>
				</div>
			</div>

			<div className="hidden md:block md:w-1/3 lg:w-1/2 bg-white relative">
				<AnimatePresence mode="wait">
					<motion.div
						key={steps[currentStep].image}
						className="absolute inset-0 bg-gray-50"
						initial={{ opacity: 0 }}
						animate={{ opacity: 1 }}
						exit={{ opacity: 0 }}
						transition={{ duration: 0.2 }}
					/>
					<motion.img
						key={`img-${steps[currentStep].image}`}
						src={steps[currentStep].image}
						alt="Dashboard Preview"
						className="absolute top-0 left-0 w-auto h-full object-cover object-left pl-8 pt-6 pb-6"
						style={{
							userSelect: "none",
							WebkitUserSelect: "none",
							MozUserSelect: "none",
							msUserSelect: "none",
						}}
						draggable={false}
						onDragStart={(e) => e.preventDefault()}
						initial={{ opacity: 0, x: 20 }}
						animate={{ opacity: 1, x: 0 }}
						exit={{ opacity: 0, x: -20 }}
						transition={{
							type: "spring",
							stiffness: 400,
							damping: 30,
							mass: 0.8,
						}}
					/>
				</AnimatePresence>
			</div>
		</div>
	);
};

export default MultiStepOnboarding;
