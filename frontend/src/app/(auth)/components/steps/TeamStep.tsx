// src/components/steps/TeamStep.tsx
"use client";

import { useOnboardingStore, useTeam } from "@/app/_store/onboarding";
import { type TeamFormData, teamFormSchema } from "@/app/_types/onboarding";
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

interface TeamStepProps {
  onNext: () => void;
}


export default function TeamStep({ onNext }: TeamStepProps) {
  const team = useTeam();
  const setTeam = useOnboardingStore((state) => state.setTeam);
  const [duplicateError, setDuplicateError] = useState<{
    [key: number]: boolean;
  }>({});

  const form = useForm<TeamFormData>({
    resolver: zodResolver(teamFormSchema),
    defaultValues: {
      members: team.length > 0
        ? team.map(email => ({ email }))
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
      const isValid = await form.trigger(`members.${index}.email`);
      const isDuplicate = checkDuplicate(currentValue, index);

      if (isValid && !isDuplicate) {
        if (index === fields.length - 1 && fields.length < 5) {
          append({ email: "" });
          setTimeout(() => {
            const inputs = document.querySelectorAll('input[type="email"]');
            const nextInput = inputs[index + 1] as HTMLInputElement;
            if (nextInput) nextInput.focus();
          }, 0);
        } else {
          const inputs = document.querySelectorAll('input[type="email"]');
          const nextInput = inputs[index + 1] as HTMLInputElement;
          if (nextInput) nextInput.focus();
        }
      }
    }
  };

  const onSubmit = async (data: TeamFormData) => {
    try {
      // Update Zustand store with just the email addresses
      setTeam(data.members.map(member => member.email));
      onNext();
    } catch (error) {
      console.error("Submission error:", error);
    }
  };

  const hasFormErrors = () => {
    const hasZodErrors = Object.keys(form.formState.errors).length > 0;
    const hasDuplicates = Object.values(duplicateError).some(Boolean);
    const hasEmptyFields = form
      .getValues()
      .members.some((member) => !member?.email?.trim());

    return hasZodErrors || hasDuplicates || hasEmptyFields;
  };

  return (
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
                              await form.trigger(`members.${index}.email`);
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
                            : form.formState.errors.members?.[index]?.email
                              ?.message}
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
          {fields.length !== 5 && (
            <Button
              type="button"
              variant="secondary"
              className="w-full text-sm border-gray-200"
              onClick={onNext}
            >
              Skip
            </Button>
          )}
        </div>
      </form>
    </Form>
  );
}
