"use client";

import {
	type CompetitorFormData,
	competitorFormSchema,
} from "@/app/_types/onboarding";
import { Button } from "@/components/ui/button";
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { Globe, Plus, X } from "lucide-react";
import Link from "next/link";
import { useMemo, useState } from "react";
import { type SubmitHandler, useFieldArray, useForm } from "react-hook-form";

const normalizeUrl = (url: string): string => {
	if (!url) return "";

	try {
		const urlObj = new URL(url.startsWith("http") ? url : `https://${url}`);
		return urlObj.hostname.toLowerCase();
	} catch {
		return url.toLowerCase();
	}
};

interface CompetitorStepProps {
	formData: CompetitorFormData;
	// biome-ignore lint/suspicious/noExplicitAny: <explanation>
	setFormData: (data: any) => void;
	onNext: () => void;
}

export default function CompetitorStep({
	formData,
	setFormData,
	onNext,
}: CompetitorStepProps) {
	const [urlErrors, setUrlErrors] = useState<{ [key: number]: boolean }>({});

	const form = useForm<CompetitorFormData>({
		resolver: zodResolver(competitorFormSchema),
		defaultValues: {
			competitors:
				formData.competitors.length > 0
					? formData.competitors
					: [{ url: "", favicon: "" }],
		},
		mode: "onBlur",
	});

	const { fields, append, remove } = useFieldArray({
		control: form.control,
		name: "competitors",
	});

	const isDuplicateUrl = (url: string, currentIndex: number): boolean => {
		const normalizedUrl = normalizeUrl(url);
		return form
			.getValues()
			.competitors.some(
				(competitor, index) =>
					index !== currentIndex &&
					normalizeUrl(competitor.url) === normalizedUrl,
			);
	};

	const hasInvalidUrls = useMemo(() => {
		return (
			Object.values(urlErrors).some((error) => error) ||
			Object.keys(urlErrors).length < fields.length
		);
	}, [urlErrors, fields.length]);

	const isValidUrl = (urlString: string): boolean => {
		// console.log("Validating URL:", urlString);
		// Define allowed protocols
		const ALLOWED_PROTOCOLS = ["http:", "https:"];

		// Regular expressions for validation
		const IP_REGEX = /^(\d{1,3}\.){3}\d{1,3}$/;
		const LOCALHOST_REGEX = /^localhost(:\d+)?$/;
		const DOMAIN_REGEX = /^([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$/;

		try {
			// Try parsing the URL
			const url = new URL(
				urlString.includes("://") ? urlString : `https://${urlString}`,
			);

			// Check protocol
			if (!ALLOWED_PROTOCOLS.includes(url.protocol)) {
				return false;
			}

			// Remove port number if present for hostname validation
			const hostname = url.hostname.includes(":")
				? url.hostname.split(":")[0]
				: url.hostname;

			// Handle special cases
			if (hostname === "localhost" || LOCALHOST_REGEX.test(hostname)) {
				return true;
			}
			if (IP_REGEX.test(hostname)) {
				// Validate IP address ranges
				const parts = hostname.split(".").map(Number);
				return parts.every((part) => part >= 0 && part <= 255);
			}

			// Validate domain name
			if (!DOMAIN_REGEX.test(hostname)) {
				return false;
			}

			// Additional checks for suspicious patterns
			const suspiciousPatterns = [
				/[^\x20-\x7E]/, // Non-printable ASCII characters
				/\s/, // Whitespace
				/[<>{}|\^~\[\]`]/, // Dangerous characters
				/javascript:/i, // JavaScript protocol
				/data:/i, // Data protocol
			];

			if (suspiciousPatterns.some((pattern) => pattern.test(urlString))) {
				return false;
			}

			return true;
		} catch {
			return false;
		}
	};

	const fetchFavicon = async (url: string, index: number) => {
		if (!isValidUrl(url)) {
			setUrlErrors((prev) => ({ ...prev, [index]: true }));
			form.setValue(`competitors.${index}.favicon`, "");
			return;
		}

		if (isDuplicateUrl(url, index)) {
			setUrlErrors((prev) => ({ ...prev, [index]: true }));
			form.setError(`competitors.${index}.url`, {
				type: "manual",
				message: "This website has already been added",
			});
			return;
		}

		try {
			const domain = new URL(url.startsWith("http") ? url : `https://${url}`)
				.hostname;
			const faviconUrl = `https://www.google.com/s2/favicons?domain=${domain}&sz=32`;

			const img = new Image();
			img.src = faviconUrl;

			await new Promise((resolve, reject) => {
				img.onload = resolve;
				img.onerror = reject;
			});

			setUrlErrors((prev) => ({ ...prev, [index]: false }));
			form.setValue(`competitors.${index}.favicon`, faviconUrl);
		} catch (error) {
			console.error("Error fetching favicon:", error);
			setUrlErrors((prev) => ({ ...prev, [index]: true }));
			form.setValue(`competitors.${index}.favicon`, "");
		}
	};

	const handleKeyPress = async (
		event: React.KeyboardEvent<HTMLInputElement>,
		index: number,
	) => {
		if (event.key === "Enter") {
			event.preventDefault();

			const currentValue = form.getValues(`competitors.${index}.url`);

			if (isDuplicateUrl(currentValue, index)) {
				form.setError(`competitors.${index}.url`, {
					type: "manual",
					message: "This website has already been added",
				});
				setUrlErrors((prev) => ({ ...prev, [index]: true }));
				return;
			}

			if (isValidUrl(currentValue)) {
				await fetchFavicon(currentValue, index);

				if (
					index === fields.length - 1 &&
					fields.length < 5 &&
					!hasInvalidUrls
				) {
					append({ url: "" });
				}

				const nextInput = document.querySelector(
					`input[name="competitors.${index + 1}.url"]`,
				) as HTMLInputElement;
				if (nextInput) {
					nextInput.focus();
				}
			}
		}
	};

	const handleRemove = (index: number) => {
		remove(index);
		// Update URL State
		setUrlErrors((prev) => {
			const newState = { ...prev };
			delete newState[index];
			for (const key of Object.keys(newState)) {
				const keyNum = Number.parseInt(key);
				if (keyNum > index) {
					newState[keyNum - 1] = newState[keyNum];
					delete newState[keyNum];
				}
			}
			return newState;
		});
	};

	const onSubmit: SubmitHandler<CompetitorFormData> = async (data) => {
		try {
			// console.log("Submitted data:", data);
			setFormData({
				...formData,
				competitors: data.competitors,
			});
			onNext();
		} catch (error) {
			console.error("Submission error:", error);
		}
	};

	return (
		<Form {...form}>
			<form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
				<div className="space-y-4">
					{fields.map((field, index) => (
						<FormField
							key={field.id}
							control={form.control}
							name={`competitors.${index}.url`}
							render={({ field }) => (
								<FormItem>
									<FormControl>
										<div className="relative flex items-center">
											{form.getValues(`competitors.${index}.favicon`) ? (
												<img
													src={form.getValues(`competitors.${index}.favicon`)}
													alt="Favicon"
													className="absolute left-3 w-4 h-4"
												/>
											) : (
												<Globe className="absolute left-3 w-4 h-4 text-gray-400" />
											)}
											<Input
												{...field}
												placeholder="Enter competitor's website"
												className={cn(
													"h-12 text-base rounded-md border-gray-200 pl-10",
													"focus:ring-2 focus:ring-offset-0 focus:ring-blue-500/20 focus:border-blue-500",
													"transition-colors duration-200",
													urlErrors[index] && "border-red-500",
												)}
												onKeyDown={(e) => handleKeyPress(e, index)}
												onBlur={async (e) => {
													field.onBlur();
													if (e.target.value) {
														await fetchFavicon(e.target.value, index);
													}
													form.trigger(`competitors.${index}.url`);
												}}
											/>
											{index > 0 && (
												<Button
													type="button"
													variant="ghost"
													size="sm"
													className="absolute right-2"
													onClick={() => handleRemove(index)}
												>
													<X className="h-4 w-4" />
												</Button>
											)}
										</div>
									</FormControl>
									<FormMessage className="text-xs text-red-500 mt-1" />
								</FormItem>
							)}
						/>
					))}
				</div>

				{fields.length < 5 && (
					<Button
						type="button"
						variant="outline"
						className="w-full h-12 border-dashed"
						onClick={() => append({ url: "", favicon: "" })}
						disabled={Object.values(urlErrors).some((error) => error)}
					>
						<Plus className="h-4 w-4 mr-2" />
						Add Competitor
					</Button>
				)}

				<div className="space-y-6">
					<Button
						type="submit"
						className={cn(
							"w-full h-12 text-base font-semibold",
							"bg-[#14171F] hover:bg-[#14171F]/90",
							"transition-colors duration-200",
							"disabled:opacity-50 disabled:cursor-not-allowed",
						)}
						size="lg"
						disabled={form.formState.isSubmitting}
					>
						{form.formState.isSubmitting ? (
							<span className="flex items-center justify-center">
								<svg
									className="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
									xmlns="http://www.w3.org/2000/svg"
									fill="none"
									viewBox="0 0 24 24"
								>
									<title>Loading</title>
									<circle
										className="opacity-25"
										cx="12"
										cy="12"
										r="10"
										stroke="currentColor"
										strokeWidth="4"
									/>
									<path
										className="opacity-75"
										fill="currentColor"
										d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
									/>
								</svg>
								Submitting...
							</span>
						) : (
							"Continue"
						)}
					</Button>

					{fields.length === 5 && (
						<p className="text-sm text-muted-foreground text-center">
							You can always add more later
						</p>
					)}

					{fields.length === 1 && (
						<p className="text-center text-sm text-gray-600">
							Already have an account?{" "}
							<Link
								href="/login"
								className="font-semibold text-gray-900 hover:underline"
							>
								Log in
							</Link>
						</p>
					)}
				</div>
			</form>
		</Form>
	);
}
