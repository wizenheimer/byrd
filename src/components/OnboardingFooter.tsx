import { cn } from "@/lib/utils";
import { Button } from "./ui/button";
import Link from "next/link";

interface OnboardingFooterProps {
    onSubmit?: () => void;
    isSubmitting?: boolean;
    buttonText?: string;
    helpText?: string;
    showLoginLink?: boolean;
    className?: string;
    disabled?: boolean;
}

const OnboardingFooter: React.FC<OnboardingFooterProps> = ({
    onSubmit,
    isSubmitting = false,
    buttonText = "Continue",
    helpText,
    showLoginLink = false,
    className,
    disabled = false
}) => {
    return (
        <div className={cn("space-y-6", className)}>
            <Button
                type="submit"
                onClick={onSubmit}
                className={cn(
                    "w-full h-12 text-base font-semibold",
                    "bg-[#14171F] hover:bg-[#14171F]/90",
                    "transition-colors duration-200",
                    "disabled:opacity-50 disabled:cursor-not-allowed",
                )}
                size="lg"
                disabled={isSubmitting || disabled}
            >
                {isSubmitting ? (
                    <span className="flex items-center justify-center">
                        {/* <LoadingSpinner /> */}
                        Submitting...
                    </span>
                ) : (
                    buttonText
                )}
            </Button>

            {helpText && (
                <p className="text-sm text-muted-foreground text-center">
                    {helpText}
                </p>
            )}

            {showLoginLink && (
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
    );
};

export type { OnboardingFooterProps };
export { OnboardingFooter };