// src/components/steps/FeaturesStep.tsx
"use client";

import { INITIAL_FEATURES } from "@/app/_constants/onboarding";
import { useEnabledFeatures, useOnboardingActions } from "@/app/_store/onboarding";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { cn } from "@/lib/utils";

interface FeaturesStepProps {
  onNext: () => void;
}

export default function FeaturesStep({ onNext }: FeaturesStepProps) {
  const enabledFeatures = useEnabledFeatures();
  const { setEnabledFeatures } = useOnboardingActions();

  const toggleFeature = (id: string) => {
    const updatedFeatures = enabledFeatures.includes(id)
      ? enabledFeatures.filter(featureId => featureId !== id)
      : [...enabledFeatures, id];

    setEnabledFeatures(updatedFeatures);
  };

  return (
    <div className="space-y-6">
      {INITIAL_FEATURES.map((feature) => (
        <div key={feature.id} className="flex items-center space-x-4">
          <Switch
            id={feature.id}
            checked={enabledFeatures.includes(feature.id)}
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
      ))}

      <div className="space-y-6">
        <Button
          className={cn(
            "w-full h-12 text-base font-semibold",
            "bg-[#14171F] hover:bg-[#14171F]/90",
            "transition-colors duration-200",
            "disabled:opacity-50 disabled:cursor-not-allowed",
          )}
          size="lg"
          onClick={onNext}
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
