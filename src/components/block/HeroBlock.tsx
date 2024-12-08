// src/components/block/HeroBlock.tsx
import { Button } from "@/components/ui/button";
import React from "react";

type HeroButtonProps = {
  label: string;
  variant?: "default" | "outline";
  href?: string;
  onClick?: () => void;
};

type HeroBadgeProps = {
  text: string;
  dotColor?: string;
};

type HeroSectionProps = {
  title: {
    desktop: string;
    mobile: string;
  };
  description: string;
  badge?: HeroBadgeProps;
  primaryButton: HeroButtonProps;
  secondaryButton?: HeroButtonProps;
  className?: string;
};

const MultiLineText = ({ text }: { text: string }) => (
  <>
    {text.split("\n").map((line, i, arr) => (
      // biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
      <React.Fragment key={i}>
        {line}
        {i < arr.length - 1 && <br />}
      </React.Fragment>
    ))}
  </>
);

const HeroSection = ({
  title,
  description,
  badge,
  primaryButton,
  secondaryButton,
  className = "",
}: HeroSectionProps) => {
  return (
    <div
      className={`w-full bg-background relative overflow-hidden ${className}`}
    >
      <div className="max-w-7xl mx-auto px-4 pt-24 pb-16 md:pb-32 lg:pb-48">
        {/* Optional Beta Badge */}
        {badge && (
          <div className="flex justify-center mb-8">
            <div className="bg-black/10 rounded-full px-3 py-1 inline-flex items-center gap-2">
              <div
                className="w-2 h-2 rounded-full animate-[statusBlink_3s_ease-in-out_infinite]"
                style={{ backgroundColor: badge.dotColor || "#22C55E" }}
              />
              <span className="text-sm">{badge.text}</span>
            </div>
          </div>
        )}

        {/* Hero Text Content */}
        <div className="text-center max-w-3xl mx-auto">
          <h1 className="text-5xl font-bold tracking-tight mb-6">
            <span className="hidden md:inline">
              <MultiLineText text={title.desktop} />
            </span>
            <span className="md:hidden">{title.mobile}</span>
          </h1>
          <p className="text-lg text-gray-600 mb-8">
            <MultiLineText text={description} />
          </p>
          <div className="flex gap-4 justify-center">
            <Button
              onClick={primaryButton.onClick}
              className="bg-black text-white hover:bg-black/90 px-8"
              asChild={!!primaryButton.href}
            >
              {primaryButton.href ? (
                <a href={primaryButton.href}>{primaryButton.label}</a>
              ) : (
                primaryButton.label
              )}
            </Button>
            {secondaryButton && (
              <Button
                variant="outline"
                className="border-gray-200"
                onClick={secondaryButton.onClick}
                asChild={!!secondaryButton.href}
              >
                {secondaryButton.href ? (
                  <a href={secondaryButton.href}>{secondaryButton.label}</a>
                ) : (
                  secondaryButton.label
                )}
              </Button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default HeroSection;
