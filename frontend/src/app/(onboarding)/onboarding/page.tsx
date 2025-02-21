"use client";

import {
	FEATURE_CARD,
	STEPS,
	STEP_INFO,
} from "@/app/(onboarding)/_constants/onboarding";
import { useOnboardingStore } from "@/app/store/onboarding";
import { Button } from "@/components/ui/button";
import { ToastAction } from "@/components/ui/toast";
import { useToast } from "@/hooks/use-toast";
import { AnimatePresence, motion } from "framer-motion";
import { ChevronLeft } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useRef } from "react";
import AuthStep from "../_components/steps/AuthStep";
import CompetitorStep from "../_components/steps/CompetitorStep";
import FeatureStep from "../_components/steps/FeaturesStep";
import ProfileStep from "../_components/steps/ProfilesStep";

const MultiStepOnboarding = () => {
	const currentStep = useOnboardingStore.use.currentStep();
	const setCurrentStep = useOnboardingStore.use.setCurrentStep();
	const resetState = useOnboardingStore.use.reset();
	const router = useRouter();
	const { toast } = useToast();
	const hasShownToast = useRef(false);

	useEffect(() => {
		// Only show toast once when component mounts and there's saved progress
		if (!hasShownToast.current && currentStep > STEPS.COMPETITOR) {
			toast({
				title: "Welcome back!",
				description: "You can pick up right where you left off.",
				duration: 5000,
				action: (
					<ToastAction
						altText="Start over"
						onClick={() => {
							resetState();
							router.push("/");
						}}
					>
						Start over
					</ToastAction>
				),
			});
			hasShownToast.current = true;
		}
	}, [currentStep, toast, resetState, router]);

	// Handle next step
	const handleNext = () => {
		if (currentStep < STEPS.AUTH) {
			if (currentStep === STEPS.COMPETITOR) {
				hasShownToast.current = true;
			}
			setCurrentStep(currentStep + 1);
		}
	};

	// Handle back step
	const handleBack = () => {
		if (currentStep > STEPS.COMPETITOR) {
			setCurrentStep(currentStep - 1);
		}
	};

	// Render step based on current step
	const renderStep = () => {
		switch (currentStep) {
			case STEPS.COMPETITOR:
				return <CompetitorStep onNext={handleNext} />;
			case STEPS.PROFILE:
				return <ProfileStep onNext={handleNext} />;
			case STEPS.FEATURES:
				return <FeatureStep cards={FEATURE_CARD} onNext={handleNext} />;
			case STEPS.AUTH:
				// Add authentication step if unauthenticated
				return <AuthStep />;
			default:
				return null;
		}
	};

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
						<span className="text-xl font-semibold">
							<Link href="/">byrd</Link>
						</span>
					</div>
					<div className="text-sm text-muted-foreground">
						Step {currentStep} of {STEPS.AUTH}
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
									{STEP_INFO[currentStep as keyof typeof STEP_INFO].title}
								</h1>
								<p className="text-base text-muted-foreground">
									{STEP_INFO[currentStep as keyof typeof STEP_INFO].description}
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
						key={STEP_INFO[currentStep as keyof typeof STEP_INFO].image}
						className="absolute inset-0 bg-gray-50"
						initial={{ opacity: 0 }}
						animate={{ opacity: 1 }}
						exit={{ opacity: 0 }}
						transition={{ duration: 0.2 }}
					/>
					<motion.img
						key={`img-${STEP_INFO[currentStep as keyof typeof STEP_INFO].image}`}
						src={STEP_INFO[currentStep as keyof typeof STEP_INFO].image}
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
