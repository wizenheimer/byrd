"use client";

// import { OnboardingHeader } from "@/components/OnboardingHeader";
// import { OnboardingLayout } from "@/components/OnboardingLayout";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { cn } from "@/lib/utils";

interface FeatureItem {
  id: string
  title: string
  description: string
  enabled: boolean
}

interface FeaturesStepProps {
  formData: {
    features: FeatureItem[];
    // biome-ignore lint/suspicious/noExplicitAny: <explanation>
    [key: string]: any;
  };
  // biome-ignore lint/suspicious/noExplicitAny: <explanation>
  setFormData: (data: any) => void;
  onNext: () => void;
  onBack: () => void;
}

export default function FeaturesStep({ formData, setFormData, onNext, onBack }: FeaturesStepProps) {

  const toggleFeature = (id: string) => {
    const updatedFeatures = formData.features.map(feature =>
      feature.id === id ? { ...feature, enabled: !feature.enabled } : feature
    );

    setFormData({
      ...formData,
      features: updatedFeatures
    });
  };

  const handleContinue = () => {
    // Validate if needed
    onNext();
  };

  return (
    <div className="space-y-6">
      {
        formData.features.map((feature) => (
          <div key={feature.id} className="flex items-center space-x-4">
            <Switch
              id={feature.id}
              checked={feature.enabled}
              onCheckedChange={() => toggleFeature(feature.id)}
              className="data-[state=checked]:bg-blue-600"
            />
            <div className="flex-1 space-y-1">
              <label
                htmlFor={feature.id}
                className="text-base font-semibold leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
              >
                {feature.title}
              </label>
              <p className="text-sm text-muted-foreground">
                {feature.description}
              </p>
            </div>
          </div>
        ))
      }

      <div className="space-y-6">
        <Button
          className={cn(
            "w-full h-12 text-base font-semibold",
            "bg-[#14171F] hover:bg-[#14171F]/90",
            "transition-colors duration-200",
            "disabled:opacity-50 disabled:cursor-not-allowed",
          )}
          size="lg"
          onClick={handleContinue}
        >
          Continue
        </Button>

        <p className="text-sm text-muted-foreground text-center">
          You can always customize them later
        </p>
      </div>
    </div>
  );
}
