"use client";

import { motion } from "framer-motion";
import { ArrowUpRight } from "lucide-react";
import { BaseLayout } from "../_components/layouts/BaseLayout";

export default function WaitlistScreen() {
	return (
		<BaseLayout header={<div className="flex items-center space-x-4" />}>
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
