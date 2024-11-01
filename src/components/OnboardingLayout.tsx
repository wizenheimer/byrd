import { cn } from "@/lib/utils";
import OnboardingPreviewPane from "./block/OnboardingPreviewPane";


interface OnboardingLayoutProps {
    children: React.ReactNode;
    previewImage: string;
    previewAlt: string;
    className?: string;
}

const OnboardingLayout: React.FC<OnboardingLayoutProps> = ({
    children,
    previewImage,
    previewAlt,
    className
}) => {
    return (
        <div className={cn(
            "flex min-h-screen flex-col lg:flex-row",
            className
        )}>
            <div className="flex flex-1 flex-col bg-white p-8 lg:p-12">
                <div className="mb-16">
                    <span className="text-xl font-semibold">byrd</span>
                </div>
                <div className="mx-auto w-full max-w-[440px] space-y-12">
                    {children}
                </div>
            </div>
            <OnboardingPreviewPane
                imageSrc={previewImage}
                altText={previewAlt}
            />
        </div>
    );
};

export type { OnboardingLayoutProps };
export { OnboardingLayout };