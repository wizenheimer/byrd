"use client";
import Footer from "@/components/Footer";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { useEffect, useState } from "react";

const errorMessages = [
	{
		message:
			"The page you're looking for cannot be found. Let's get you back on track.",
	},
	{
		message:
			"Looks like this page has moved or no longer exists. We can help you find what you need.",
	},
	{
		message:
			"This link appears to be broken, but we know just where to redirect you.",
	},
	{
		message:
			"Sorry about that! The page you requested isn't here, but we know where you can find what you need.",
	},
	{
		message:
			"We can't seem to find the page you're looking for. Let's help you find the right place.",
	},
	{
		message:
			"This page seems to be missing, but don't worry - we'll help you get to where you need to go.",
	},
	{
		message:
			"The requested page couldn't be found. Let us point you in the right direction.",
	},
	{
		message:
			"We couldn't find what you're looking for, but we're here to help you find the right path.",
	},
	{
		message:
			"Something's not quite right with this link. Let's get you to the right place.",
	},
	{
		message:
			"This page appears to be missing, but we'll help you find what you're looking for.",
	},
];

export default function NotFound() {
	// Start with a default message
	const [message, setMessage] = useState(errorMessages[0].message);

	// Only update the message on the client side after mount
	useEffect(() => {
		const randomMessage =
			errorMessages[Math.floor(Math.random() * errorMessages.length)];
		setMessage(randomMessage.message);
	}, []);

	return (
		<>
			<div className="w-full bg-background relative overflow-hidden">
				<div className="max-w-7xl mx-auto px-4 py-24">
					<div className="text-center max-w-3xl mx-auto">
						<h1 className="text-5xl font-bold tracking-tight mb-6">
							Page Not Found
						</h1>
						<p className="text-lg text-gray-600 mb-8">{message}</p>
						<div className="flex gap-4 justify-center">
							<Link href="/">
								<Button className="bg-black text-white hover:bg-black/90 px-8">
									Back to Homepage
								</Button>
							</Link>
							<Link href="/help">
								<Button variant="outline" className="border-gray-200">
									Help Center
								</Button>
							</Link>
						</div>
					</div>
				</div>
			</div>
			<Footer />
		</>
	);
}
