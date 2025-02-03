"use client";

import LoadingStep from "@/app/(onboarding)/components/steps/LoadingStep";
import { UserButton, useUser } from "@clerk/nextjs";
import { motion } from "framer-motion";
import { ArrowUpRight, ExternalLink, Plus } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { BaseLayout } from "../components/layouts/BaseLayout";

export default function WaitlistScreen() {
	const { isLoaded, isSignedIn } = useUser();
	const router = useRouter();

	useEffect(() => {
		if (isLoaded && !isSignedIn) {
			router.push("/create");
		}
	}, [isLoaded, isSignedIn, router]);

	const handleSkipWaitlist = () => {
		window.location.href = "https://cal.com/nayann/byrd";
	};

	const handleCreateWorkspace = () => {
		router.push("/create");
	};

	if (!isLoaded || !isSignedIn) {
		return <LoadingStep />;
	}

	return (
		<BaseLayout
			header={
				<div className="flex items-center space-x-4">
					<UserButton
						appearance={{
							elements: {
								// This removes the default menu items
								userButtonPopoverActionButton: "hidden",
								userButtonPopoverActionButtonText: "hidden",
								userButtonPopoverFooter: "hidden",
							},
						}}
					>
						<UserButton.MenuItems>
							<UserButton.Action
								label="Create workspace"
								labelIcon={<Plus className="w-4 h-4" />}
								onClick={handleCreateWorkspace}
							/>
							<UserButton.Action
								label="Skip waitlist"
								labelIcon={<ExternalLink className="w-4 h-4" />}
								onClick={handleSkipWaitlist}
							/>
						</UserButton.MenuItems>
					</UserButton>
				</div>
			}
		>
			<motion.div
				initial={{ opacity: 0, y: 20 }}
				animate={{ opacity: 1, y: 0 }}
				transition={{ duration: 0.5 }}
				className="space-y-4"
			>
				<div className="flex items-center mb-4">
					<h1 className="text-6xl font-bold text-black">Stay tuned</h1>
					<ArrowUpRight className="w-12 h-12 ml-2 text-black" />
				</div>
				<p className="text-md text-gray-600 max-w-md">
					{"You're on the list! We'll reach out soon."}
				</p>
			</motion.div>
		</BaseLayout>
	);
}
