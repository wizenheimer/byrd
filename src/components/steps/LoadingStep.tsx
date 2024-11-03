"use client";

import { AnimatePresence, motion } from "framer-motion";
import Link from "next/link";
import { useEffect, useState } from "react";

const loadingMessages = [
	"Preparing your experience",
	"To load, or not to load, that is the question",
	"Loading like nobody's watching",
	"This loading screen is powered by AI",
	"Leveraging blockchain for better loading times",
	"Loading goes brrrrrrr",
	"Making loading 10x better (YC W24)",
	"This loading screen is our moat",
	"Doing things that don't scale (like this loading bar)",
	"PM said this loading screen increases user engagement by 300%",
	"This loading increased our north star metric somehow",
	"POV: You're watching a loading screen",
	"Loading.exe has stopped loading.exe",
	"Loading screen: 🗿",
	"POV: You're watching paint.exe dry 🗿",
	"Yo dawg, I heard you like loading screens",
	"Loading screen got that dawg in it fr fr",
	"POV: You're in a staring contest with loading animation",
	"Making things perfect",
	"Loading screen (YC W24) - 'Stripe for waiting'",
	"Loading our 'Why YC rejected us' Medium post",
	"Loading things that don't scale",
	"Make Loading People Want",
	"Almost there",
	"Loading awesome content",
	"Gathering the good stuff",
	"Just adding final touches",
	"This is worth the wait",
	"Brewing something special",
	"Making magic happen",
	"Double-checking everything",
	"To load, or not to load, that is the question",
	"Loading like nobody's watching",
	"This loading screen is powered by AI",
	"Leveraging blockchain for better loading times",
	"Loading goes brrrrrrr",
	"Making loading 10x better (YC W24)",
	"This loading screen is our moat",
	"Doing things that don't scale (like this loading bar)",
	"PM said this loading screen increases user engagement by 300%",
	"This loading increased our north star metric somehow",
	"POV: You're watching a loading screen",
	"Loading.exe has stopped loading.exe",
	"Loading screen: 🗿",
	"POV: You're watching paint.exe dry 🗿",
	"Yo dawg, I heard you like loading screens",
	"Loading screen got that dawg in it fr fr",
	"POV: You're in a staring contest with loading animation",
];

const LoadingStep = () => {
	const [currentMessage, setCurrentMessage] = useState(
		loadingMessages[Math.floor(Math.random() * loadingMessages.length)],
	);

	// Change message every 3 seconds
	useEffect(() => {
		const messageInterval = setInterval(() => {
			setCurrentMessage(
				loadingMessages[Math.floor(Math.random() * loadingMessages.length)],
			);
		}, 3000);

		return () => clearInterval(messageInterval);
	}, []);

	return (
		<motion.div
			initial={{ opacity: 0 }}
			animate={{ opacity: 1 }}
			exit={{ opacity: 0 }}
			transition={{ duration: 0.5 }}
			className="min-h-screen bg-white text-black flex flex-col"
		>
			<nav>
				<div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
					<Link
						href="/"
						className="text-xl font-bold hover:scale-105 transition-transform"
					>
						byrd
					</Link>
				</div>
			</nav>

			<main className="flex-grow flex flex-col justify-center items-center px-4 text-center">
				<motion.h1
					initial={{ y: -20, opacity: 0 }}
					animate={{ y: 0, opacity: 1 }}
					transition={{ delay: 0.2, duration: 0.5 }}
					className="text-6xl font-bold text-black mb-6"
				>
					Just a moment
				</motion.h1>

				<AnimatePresence mode="wait">
					<motion.p
						key={currentMessage}
						initial={{ y: 20, opacity: 0 }}
						animate={{ y: 0, opacity: 1 }}
						exit={{ y: -20, opacity: 0 }}
						transition={{ duration: 0.3 }}
						className="text-lg text-gray-600 max-w-md mb-8"
					>
						{currentMessage}
					</motion.p>
				</AnimatePresence>

				<motion.div
					initial={{ scale: 0 }}
					animate={{ scale: 1 }}
					transition={{
						delay: 0.5,
						type: "spring",
						stiffness: 200,
						damping: 20,
					}}
				>
					<div className="w-16 h-16 border-4 border-black border-t-transparent rounded-full animate-spin" />
				</motion.div>
			</main>

			<motion.footer
				initial={{ y: 20, opacity: 0 }}
				animate={{ y: 0, opacity: 1 }}
				transition={{ delay: 0.7, duration: 0.5 }}
				className="py-8 px-6"
			>
				<div className="max-w-7xl mx-auto">
					<div className="flex justify-between items-center text-sm text-gray-600">
						<p>© 2024 byrd. All rights reserved.</p>
						<div className="flex space-x-4">
							<Link
								href="/terms"
								className="hover:text-black transition-colors"
							>
								Terms
							</Link>
							<Link
								href="/privacy"
								className="hover:text-black transition-colors"
							>
								Privacy
							</Link>
						</div>
					</div>
				</div>
			</motion.footer>
		</motion.div>
	);
};

export default LoadingStep;
