"use client";

import { ToastAction } from "@/components/ui/toast";
import { useToast } from "@/hooks/use-toast";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useRef } from "react";
import { handleSlackCallback } from "../../_actions/onboarding";
import Loading from "../../_components/Loading";

export default function CompletePage() {
	const { toast } = useToast();
	const router = useRouter();
	const searchParams = useSearchParams();
	const installationStarted = useRef(false);

	useEffect(() => {
		const completeInstallation = async () => {
			// Prevent multiple installations
			if (installationStarted.current) return;
			installationStarted.current = true;

			const code = searchParams.get("code");
			const state = searchParams.get("state");

			if (!code || !state) {
				toast({
					variant: "destructive",
					title: "Installation Failed",
					description: "Missing required parameters",
					action: (
						<ToastAction
							altText="Try again"
							onClick={() => router.push("/onboarding")}
						>
							Try again
						</ToastAction>
					),
				});
				router.push("/onboarding");
				return;
			}

			try {
				const result = await handleSlackCallback(code, state);

				if (!result.success) {
					throw new Error(result.error);
				}

				const { deep_link } = result.data;
				if (deep_link) {
					window.location.href = deep_link;
				} else {
					throw new Error("Invalid response from server");
				}
			} catch (error) {
				toast({
					variant: "destructive",
					title: "Installation Failed",
					description:
						error instanceof Error
							? error.message
							: "Failed to complete installation",
					action: (
						<ToastAction
							altText="Try again"
							onClick={() => router.push("/onboarding")}
						>
							Try again
						</ToastAction>
					),
				});
				router.push("/onboarding");
			}
		};

		completeInstallation();
	}, [searchParams, toast, router]);

	return <Loading />;
}
