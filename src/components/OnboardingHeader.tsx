import { cn } from "@/lib/utils";

interface OnboardingHeaderProps {
    title: string;
    description: string;
    className?: string;
}

const OnboardingHeader: React.FC<OnboardingHeaderProps> = ({
    title,
    description,
    className
}) => {
    return (
        <div className={cn("space-y-3", className)}>
            <h1 className="text-3xl font-bold tracking-tight">
                {title}
            </h1>
            <p className="text-base text-muted-foreground">
                {description}
            </p>
        </div>
    );
};

export type { OnboardingHeaderProps };
export { OnboardingHeader };