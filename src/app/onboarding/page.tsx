"use client";
import { useState } from "react";
import CompetitorStep from "../../components/steps/CompetitorStep";
import FeaturesStep from "../../components/steps/FeaturesStep";
import ChannelsStep from "../../components/steps/ChannelsStep";
import TeamStep from "../../components/steps/TeamStep";
import AuthStep from "../../components/steps/AuthStep";
import { OnboardingLayout } from "@/components/OnboardingLayout";
import { OnboardingHeader } from "@/components/OnboardingHeader";
import { Inbox, Megaphone, Rss, Share2 } from "lucide-react";
import type { ChannelCard } from "../_types/onboarding";

const MultiStepOnboarding = () => {
  const [currentStep, setCurrentStep] = useState(1);
  const [formData, setFormData] = useState({
    competitors: [],
    features: [
      { id: "1", title: "Product", description: "Catch product evolution in real-time", enabled: true },
      { id: "2", title: "Pricing", description: "Never be the last to know about a price war", enabled: false },
      { id: "3", title: "Partnership", description: "Track who's teaming up with whom", enabled: false },
      { id: "4", title: "Branding", description: "Monitor messaging shifts, and identity changes", enabled: false },
      { id: "5", title: "Positioning", description: "Track narratives before they go mainstream", enabled: false },
    ],
    channels: ["inbox", "mentions"],
    team: []
  });

  interface Step {
    title: string;
    description: string;
    image: string;
  }


  const steps: Record<number, Step> = {
    1: {
      title: "Your Market, Your Rules",
      description: "Pick your targets. Add up to 5 competitors.",
      image: "/onboarding/first.png"
    },
    2: {
      title: "Measure What Matters",
      description: "Choose your signals. Cut through the noise.",
      image: "/onboarding/second.png"
    },
    3: {
      title: "Never Miss A Beat",
      description: "Your competitors are everywhere. So are we.",
      image: "/onboarding/third.png"
    },
    4: {
      title: "Build Your War Room",
      description: "Business is a team sport. Bring in your heavy hitters.",
      image: "/onboarding/four.png"
    },
    5: {
      title: "You're almost there",
      description: "Quick auth, then let's get started.",
      image: "/onboarding/five.png"
    }
  };

  const channelCards: ChannelCard[] = [
    {
      id: "inbox",
      icon: Inbox,
      title: "Inbox",
      description: "Monitor the direct line to their customers",
    },
    {
      id: "social",
      icon: Share2,
      title: "Social",
      description: "Follow their social playbook as it unfolds",
    },
    {
      id: "mentions",
      icon: Megaphone,
      title: "Mentions",
      description: "Beat them to their own announcement",
    },
    {
      id: "content",
      icon: Rss,
      title: "Content",
      description: "Watch their messaging evolve, post by post",
    },
  ];

  const handleNext = () => {
    if (currentStep < 5) {
      setCurrentStep(prev => prev + 1);
    }
  };

  const handleBack = () => {
    if (currentStep > 1) {
      setCurrentStep(prev => prev - 1);
    }
  };

  const renderStep = () => {
    switch (currentStep) {
      case 1:
        return <CompetitorStep formData={formData} setFormData={setFormData} onNext={handleNext} />;
      case 2:
        return <FeaturesStep formData={formData} setFormData={setFormData} onNext={handleNext} onBack={handleBack} />;
      case 3:
        return <ChannelsStep
          formData={formData}
          setFormData={setFormData}
          onNext={handleNext}
          onBack={handleBack}
          cards={channelCards}
        />;
      case 4:
        return <TeamStep formData={formData} setFormData={setFormData} onNext={handleNext} onBack={handleBack} />;
      case 5:
        return <AuthStep onBack={handleBack} />;
      default:
        return null;
    }
  };

  return (
    <OnboardingLayout
      previewImage={steps[currentStep].image}
      previewAlt="Dashboard Preview"
    >
      <OnboardingHeader
        title={steps[currentStep].title}
        description={steps[currentStep].description}
      />

      {renderStep()}
    </OnboardingLayout>
  );
}

export default MultiStepOnboarding
