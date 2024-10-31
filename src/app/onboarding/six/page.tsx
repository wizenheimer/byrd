"use client";

import OnboardingPreviewPane from "@/components/block/OnboardingPreviewPane";
import Link from "next/link";

export default function Component() {
	return (
		<div className="flex min-h-screen flex-col lg:flex-row">
			{/* Left Side - Waitlist */}
			<div className="flex flex-1 flex-col justify-between p-8 lg:p-12">
				<div>
					<span className="text-xl font-semibold">byrd</span>
				</div>
				<div className="mx-auto w-full max-w-[480px] space-y-4">
					<div className="space-y-4 text-center">
						<h1 className="text-3xl font-bold tracking-tight">
							Yeah, There&apos;s a Waitlist
						</h1>
						<div className="space-y-1">
							<p className="text-base">
								We&apos;d rather nail it with a few than fail it with many.
							</p>
							<p className="text-base">
								Stay tuned. Rolling out access to 7 new teams every month.
							</p>
						</div>
					</div>
				</div>
				<div className="text-center">
					<Link className="text-sm text-muted-foreground" href={"/"}>
						Can&apos;t Wait? Don&apos;t.
					</Link>
				</div>
			</div>

			{/* Right Side - Dashboard Preview */}
			<OnboardingPreviewPane
				imageSrc="/onboarding/six.png"
				altText="Dashboard Preview"
			/>
		</div>
	);
}
