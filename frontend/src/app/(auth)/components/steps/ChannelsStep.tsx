// src/components/steps/ChannelsStep.tsx
"use client";

import { useChannels, useOnboardingStore } from "@/app/_store/onboarding";
import type { ChannelCard } from "@/app/_types/onboarding";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface ChannelsStepProps {
  cards: ChannelCard[];
  onNext: () => void;
}

export default function ChannelsStep({ cards, onNext }: ChannelsStepProps) {
  const selectedChannels = useChannels();
  const setChannels = useOnboardingStore((state) => state.setChannels);

  const toggleChannel = (id: string) => {
    const newChannels = selectedChannels.includes(id)
      ? selectedChannels.filter((channelId) => channelId !== id)
      : [...selectedChannels, id];

    setChannels(newChannels);
  };

  return (
    <div className="space-y-6">
      {cards.map(({ id, icon, title, description }) => (
        <ChannelButton
          key={id}
          id={id}
          icon={icon}
          title={title}
          description={description}
          isSelected={selectedChannels.includes(id)}
          onClick={() => toggleChannel(id)}
        />
      ))}
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
  );
}

const ChannelButton = ({
  // id,
  icon: Icon,
  title,
  description,
  isSelected,
  onClick,
}: {
  id: string;
  icon: React.ComponentType<{ className?: string }>;
  title: string;
  description: string;
  isSelected: boolean;
  onClick: () => void;
}) => (
  <button
    type="button"
    onClick={onClick}
    className={cn(
      "relative flex w-full items-start gap-4 rounded-xl border-2 p-4 text-left transition-colors",
      isSelected
        ? "border-primary bg-primary/5"
        : "border-border bg-background hover:border-primary/50",
    )}
  >
    <Icon className="mt-1 size-5 shrink-0" />
    <div className="space-y-1">
      <div className="font-medium">{title}</div>
      <div className="text-sm text-muted-foreground">{description}</div>
    </div>
    {isSelected && (
      <div className="absolute right-3 top-1/2 size-4 -translate-y-1/2 rounded-full bg-primary" />
    )}
  </button>
);
