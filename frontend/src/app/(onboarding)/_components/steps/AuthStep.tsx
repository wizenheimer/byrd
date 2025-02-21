"use client";

import { Button } from "@/components/ui/button";
import { Slack } from "lucide-react";

export default function AuthStep() {
	return (
		<div className="space-y-4">
			<Button
				variant="outline"
				className="relative h-12 w-full justify-center text-base font-normal"
				onClick={() => console.log("slack auth triggered")}
			>
				<div className="absolute left-4 size-5">
					<Slack />
				</div>
				Sign in with Slack
			</Button>
		</div>
	);
}
