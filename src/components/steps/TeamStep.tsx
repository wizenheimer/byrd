"use client";

// import { OnboardingHeader } from "@/components/OnboardingHeader";
// import { OnboardingLayout } from "@/components/OnboardingLayout";
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
import Avatar from "boring-avatars";
import { Plus, X } from "lucide-react";
import { useState } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { z } from "zod";

const teamMemberSchema = z.object({
	email: z.string().email("Invalid email address").min(1, "Email is required"),
});

const teamFormSchema = z.object({
	members: z
		.array(teamMemberSchema)
		.min(1, "Add at least one team member")
		.max(5, "Maximum 5 team members allowed")
		.refine((members) => {
			const emails = members.map((m) => m.email.toLowerCase());
			return new Set(emails).size === emails.length;
		}, "Duplicate email addresses are not allowed"),
});

type TeamFormData = z.infer<typeof teamFormSchema>;

interface TeamStepProps {
	formData: {
		team: { email: string }[];
		// biome-ignore lint/suspicious/noExplicitAny: <explanation>
		[key: string]: any;
	};
	// biome-ignore lint/suspicious/noExplicitAny: <explanation>
	setFormData: (data: any) => void;
	onNext: () => void;
	onBack: () => void;
}

export default function TeamStep({ formData, setFormData, onNext, onBack }: TeamStepProps) {
	const [duplicateError, setDuplicateError] = useState<{ [key: number]: boolean }>({});

	const form = useForm<TeamFormData>({
		resolver: zodResolver(teamFormSchema),
		defaultValues: {
			members: formData.team.length > 0
				? formData.team
				: [{ email: "" }],
		},
	});

	const { fields, append, remove } = useFieldArray({
		control: form.control,
		name: "members",
	});

	const checkDuplicate = (email: string, currentIndex: number) => {
		const hasMatch = form
			.getValues()
			.members.some(
				(member, idx) =>
					idx !== currentIndex &&
					member.email.toLowerCase() === email.toLowerCase(),
			);

		setDuplicateError((prev) => ({
			...prev,
			[currentIndex]: hasMatch,
		}));

		return hasMatch;
	};

	const handleKeyDown = async (
		event: React.KeyboardEvent<HTMLInputElement>,
		index: number,
	) => {
		if (event.key === "Enter") {
			event.preventDefault();

			const currentValue = form.getValues(`members.${index}.email`);

			// Validate the current field first
			const isValid = await form.trigger(`members.${index}.email`);
			const isDuplicate = checkDuplicate(currentValue, index);

			if (isValid && !isDuplicate) {
				// If we're at the last field and haven't hit the limit, add a new one
				if (index === fields.length - 1 && fields.length < 5) {
					append({ email: "" });
					// Focus will be handled after render
					setTimeout(() => {
						const inputs = document.querySelectorAll('input[type="email"]');
						const nextInput = inputs[index + 1] as HTMLInputElement;
						if (nextInput) nextInput.focus();
					}, 0);
				} else {
					// Focus the next input if available
					const inputs = document.querySelectorAll('input[type="email"]');
					const nextInput = inputs[index + 1] as HTMLInputElement;
					if (nextInput) nextInput.focus();
				}
			}
		}
	};

	const onSubmit = async (data: TeamFormData) => {
		try {
			console.log("Submitted data:", data);
			setFormData({
				...formData,
				team: data.members
			});
			onNext();
		} catch (error) {
			console.error("Submission error:", error);
		}
	};

	const hasFormErrors = () => {
		// Check for validation errors
		const hasZodErrors = Object.keys(form.formState.errors).length > 0;

		// Check for duplicate errors
		const hasDuplicates = Object.values(duplicateError).some(Boolean);

		// Check for empty fields
		const hasEmptyFields = form
			.getValues()
			.members.some((member) => !member.email.trim());

		return hasZodErrors || hasDuplicates || hasEmptyFields;
	};

	return (
		// <OnboardingLayout
		// 	previewImage="/onboarding/four.png"
		// 	previewAlt="Dashboard Preview"
		// >
		// 	<OnboardingHeader
		// 		title="Build Your War Room"
		// 		description="Business is a team sport. Bring in your heavy hitters."
		// 	/>

		<Form {...form}>
			<form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
				<div className="space-y-4">
					{fields.map((field, index) => (
						<div key={field.id} className="group relative">
							<div className="flex items-center gap-3 rounded-lg bg-muted/50 p-3">
								<div className="flex size-10 shrink-0 items-center justify-center">
									<Avatar
										size={40}
										name={field.id}
										variant="beam"
										colors={[
											"#2463eb",
											"#4b80ee",
											"#729df1",
											"#99baf4",
											"#c0d7f7",
										]}
									/>
								</div>
								<div className="flex-1">
									<FormField
										control={form.control}
										name={`members.${index}.email`}
										render={({ field }) => (
											<FormItem>
												<FormControl>
													<Input
														{...field}
														className={cn(
															"h-10 bg-background",
															duplicateError[index] &&
															"border-red-500 focus:border-red-500",
															form.formState.errors.members?.[index] &&
															"border-red-500 focus:border-red-500",
														)}
														type="email"
														placeholder="Email address"
														onKeyDown={(e) => handleKeyDown(e, index)}
														onChange={(e) => {
															field.onChange(e);
															if (duplicateError[index]) {
																checkDuplicate(e.target.value, index);
															}
														}}
														onBlur={async (e) => {
															field.onBlur();
															await form.trigger(
																`members.${index}.email`,
															);
															if (e.target.value) {
																checkDuplicate(e.target.value, index);
															}
														}}
														onPaste={(e) => {
															e.preventDefault();
															const pastedText =
																e.clipboardData.getData("text");
															const emails = pastedText
																.split(/[\s,;]+/)
																.filter((email) => email.includes("@"))
																.slice(0, 5 - fields.length + 1);

															if (!checkDuplicate(emails[0], index)) {
																field.onChange(emails[0] || "");
																form.trigger(`members.${index}.email`);

																for (const email of emails.slice(1)) {
																	if (
																		fields.length < 5 &&
																		!form
																			.getValues()
																			.members.some(
																				(m) =>
																					m.email.toLowerCase() ===
																					email.toLowerCase(),
																			)
																	) {
																		append({ email });
																	}
																}
															}
														}}
													/>
												</FormControl>
												<FormMessage className="text-xs">
													{duplicateError[index]
														? "This email is already added"
														: form.formState.errors.members?.[index]
															?.email?.message}
												</FormMessage>
											</FormItem>
										)}
									/>
								</div>
							</div>
							{index > 0 && (
								<Button
									type="button"
									variant="ghost"
									size="sm"
									className="absolute -right-2 -top-2 size-6 rounded-full p-0 opacity-0 transition-opacity group-hover:opacity-100"
									onClick={() => {
										remove(index);
										setDuplicateError((prev) => {
											const newState = { ...prev };
											delete newState[index];
											return newState;
										});
									}}
								>
									<X className="size-4" />
								</Button>
							)}
						</div>
					))}
				</div>

				{fields.length < 5 && (
					<Button
						type="button"
						variant="outline"
						className={cn(
							"w-full h-12 border-dashed",
							"disabled:opacity-50 disabled:cursor-not-allowed",
						)}
						onClick={() => append({ email: "" })}
						disabled={hasFormErrors()}
					>
						<Plus className="h-4 w-4 mr-2" />
						Add Team Member
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
						disabled={
							form.formState.isSubmitting ||
							Object.values(duplicateError).some(Boolean)
						}
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
							You can always invite more team members later
						</p>
					)}
				</div>
			</form>
		</Form>
		// </OnboardingLayout>
	);
}
