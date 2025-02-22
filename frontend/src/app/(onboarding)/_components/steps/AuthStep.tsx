// frontend/src/app/(onboarding)/_components/steps/AuthStep.tsx
import { Button } from "@/components/ui/button";
import { Slack } from "lucide-react";

export default function AuthStep() {
  const handleSlackInstall = async () => {
    try {
      // First get the OAuth URL from your backend
      const response = await fetch('/api/oauth/slack/init', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          competitors: [],
          features: [],
          profiles: []
        })
      });

      if (!response.ok) {
        throw new Error('Failed to get OAuth URL');
      }

      const { oauth_url } = await response.json();

      // Redirect to Slack's OAuth page
      window.location.href = oauth_url;
    } catch (error) {
      console.error('Failed to initiate Slack installation:', error);
    }
  };

  return (
    <div className="space-y-4">
      <Button
        variant="outline"
        className="relative h-12 w-full justify-center text-base font-normal"
        onClick={handleSlackInstall}
      >
        <div className="absolute left-4 size-5">
          <Slack />
        </div>
        Sign in with Slack
      </Button>
    </div>
  );
}
