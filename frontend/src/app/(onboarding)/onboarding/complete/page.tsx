// app/(onboarding)/complete/page.tsx
'use client';

import { useEffect, useRef, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useToast } from "@/hooks/use-toast";
import { ToastAction } from "@/components/ui/toast";
import { handleSlackCallback } from '../../_actions/onboarding';
import Loading from '../../_components/Loading';
import { Check } from 'lucide-react';

export default function CompletePage() {
  const { toast } = useToast();
  const router = useRouter();
  const searchParams = useSearchParams();
  const installationStarted = useRef(false);
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');

  useEffect(() => {
    const completeInstallation = async () => {
      if (installationStarted.current) return;
      installationStarted.current = true;

      const code = searchParams.get('code');
      const state = searchParams.get('state');

      if (!code || !state) {
        toast({
          title: "Installation Failed",
          duration: 10000,
          description: "Missing required parameters",
          action: (
            <ToastAction altText="Try again" onClick={() => router.push('/onboarding')}>
              Try again
            </ToastAction>
          ),
        });
        return;
      }

      try {
        const result = await handleSlackCallback(code, state);

        if (!result.success) {
          throw new Error(result.error);
        }

        const { deep_link } = result.data;
        if (deep_link) {
          setStatus('success');
          // Redirect after a short delay to show success state
          setTimeout(() => {
            window.location.href = deep_link;
          }, 1000);
        } else {
          throw new Error('Invalid response from server');
        }
      } catch (error) {
        setStatus('error');
        toast({
          title: "Installation Failed",
          description: error instanceof Error ? error.message : "Failed to complete installation",
          action: (
            <ToastAction altText="Try again" onClick={() => router.push('/onboarding')}>
              Try again
            </ToastAction>
          ),
        });
        router.push('/onboarding');
      }
    };

    completeInstallation();
  }, [searchParams, toast, router]);

  if (status === 'success') {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center space-y-4">
          <div className="h-12 w-12 rounded-full bg-green-100 flex items-center justify-center mx-auto">
            <Check className="h-6 w-6 text-green-600" />
          </div>
          <h2 className="text-xl font-semibold">Installation Successful!</h2>
          <p className="text-gray-600">Redirecting you to Slack...</p>
        </div>
      </div>
    );
  }

  return <Loading />;
}
