"use client";

import { Button } from "@/components/ui/button";
import { useSignIn } from "@clerk/nextjs";
import type { OAuthStrategy } from "@clerk/types";
import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useState } from "react";

const LoadingSkeleton = () => {
	return (
		<div className="fixed inset-0 bg-white">
			<div className="flex h-full">
				{/* Left side */}
				<div className="flex-1 p-8 lg:p-12">
					{/* Logo area */}
					<div className="h-7 w-16 bg-gray-100 animate-pulse rounded mb-16" />

					{/* Content container */}
					<div className="max-w-[440px] mx-auto space-y-6">
						<div className="h-8 w-3/4 bg-gray-100 animate-pulse rounded" />
						<div className="h-5 w-2/3 bg-gray-100 animate-pulse rounded" />
						<div className="h-12 w-full bg-gray-100 animate-pulse rounded mt-6" />
					</div>
				</div>

				{/* Right side */}
				<div className="hidden md:block md:w-1/3 lg:w-1/2 bg-gray-50">
					<div className="h-full bg-gray-100 animate-pulse m-6 ml-8" />
				</div>
			</div>
		</div>
	);
};

function AuthStep() {
	const { signIn } = useSignIn();
	const [isLoading, setIsLoading] = useState(false);

	const handleOAuthSignIn = async (strategy: OAuthStrategy) => {
		if (!signIn) return;

		try {
			setIsLoading(true);
			await signIn.authenticateWithRedirect({
				strategy,
				redirectUrl: "/sso-login-callback",
				redirectUrlComplete: "/dashboard",
			});
		} catch (error) {
			console.error("OAuth error", error);
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div className="space-y-3">
			<Button
				variant="outline"
				className="relative h-12 w-full justify-center text-base font-normal"
				onClick={() => handleOAuthSignIn("oauth_google")}
				disabled={isLoading}
			>
				<div className="absolute left-4 size-5">
					<svg viewBox="0 0 24 24" className="size-5">
						<title>Google Icon</title>
						<path
							d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
							fill="#4285F4"
						/>
						<path
							d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
							fill="#34A853"
						/>
						<path
							d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
							fill="#FBBC05"
						/>
						<path
							d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
							fill="#EA4335"
						/>
					</svg>
				</div>
				Continue with Google
			</Button>
		</div>
	);
}

const LoginPage = () => {
	const { isLoaded } = useSignIn();
	const [isMounted, setIsMounted] = useState(false);

	useEffect(() => {
		setIsMounted(true);
	}, []);

	// Show skeleton until everything is ready
	if (!isMounted || !isLoaded) {
		return <LoadingSkeleton />;
	}

	return (
		<div className="flex min-h-screen flex-col lg:flex-row">
			<div className="flex flex-1 flex-col bg-white p-8 lg:p-12">
				<nav className="mb-16 flex items-center justify-between">
					<div className="flex items-center">
						<span className="text-xl font-semibold">byrd</span>
					</div>
				</nav>

				<div className="mx-auto w-full max-w-[440px] space-y-12">
					<AnimatePresence mode="wait">
						<motion.div
							key={1}
							initial={{ opacity: 0, y: 20 }}
							animate={{ opacity: 1, y: 0 }}
							exit={{ opacity: 0, y: -20 }}
							transition={{ duration: 0.3 }}
							className="space-y-6"
						>
							<div className="space-y-3">
								<h1 className="text-2xl font-bold tracking-tight">
									{"We're almost there"}
								</h1>
								<p className="text-base text-muted-foreground">
									Quick auth and we make this official
								</p>
							</div>
							<AuthStep />
						</motion.div>
					</AnimatePresence>
				</div>
			</div>

			<div className="hidden md:block md:w-1/3 lg:w-1/2 bg-white relative">
				<AnimatePresence mode="wait">
					<motion.div
						key="/onboarding/five.png"
						className="absolute inset-0 bg-gray-50"
						initial={{ opacity: 0 }}
						animate={{ opacity: 1 }}
						exit={{ opacity: 0 }}
						transition={{ duration: 0.2 }}
					/>
					<motion.img
						key={`img-${"/onboarding/five.png"}`}
						src={"/onboarding/five.png"}
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

export default LoginPage;
