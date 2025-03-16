"use client";

import { Button } from "@/components/ui/button";
import { ToastAction } from "@/components/ui/toast";
import { useToast } from "@/hooks/use-toast";
import { AnimatePresence, motion } from "framer-motion";
import { Check, Loader2 } from "lucide-react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useRef, useState } from "react";
import { handleSlackCallback } from "../../_actions/onboarding";

function CompletePageContent() {
  const { toast } = useToast();
  const router = useRouter();
  const searchParams = useSearchParams();
  const installationStarted = useRef(false);
  const [status, setStatus] = useState("loading");

  useEffect(() => {
    const completeInstallation = async () => {
      if (installationStarted.current) return;
      installationStarted.current = true;

      const code = searchParams.get("code");
      const state = searchParams.get("state");

      if (!code || !state) {
        toast({
          title: "Installation Failed",
          duration: 10000,
          description: "Missing required parameters",
          action: (
            <ToastAction
              altText="Try again"
              onClick={() => router.push("/onboarding")}
            >
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
          setStatus("success");
          setTimeout(() => {
            window.location.href = deep_link;
          }, 1500);
        } else {
          throw new Error("Invalid response from server");
        }
      } catch (error) {
        setStatus("error");
        toast({
          title: "Installation Failed",
          description:
            error instanceof Error
              ? error.message
              : "Failed to complete installation",
          action: (
            <ToastAction
              altText="Try again"
              onClick={() => router.push("/onboarding")}
            >
              Try again
            </ToastAction>
          ),
        });
        router.push("/onboarding");
      }
    };

    completeInstallation();
  }, [searchParams, toast, router]);

  return (
    <div className="flex min-h-screen flex-col lg:flex-row">
      <div className="flex flex-1 flex-col bg-white p-8 lg:p-12">
        <nav className="mb-16 flex items-center justify-between">
          <Link href="/" className="text-xl font-semibold">
            byrd
          </Link>
        </nav>

        <div className="mx-auto w-full max-w-[440px] space-y-12">
          <AnimatePresence mode="wait">
            <motion.div
              key={status}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
              className="space-y-6"
            >
              <div className="space-y-3">
                <h1 className="text-2xl font-bold tracking-tight">
                  {status === "loading"
                    ? "Installing Slack"
                    : "Installation Complete"}
                </h1>
                <p className="text-base text-muted-foreground">
                  {status === "loading"
                    ? "Setting up your workspace..."
                    : "Successfully connected with Slack"}
                </p>
              </div>

              <div className="space-y-4">
                <Button
                  variant="outline"
                  className="relative h-12 w-full justify-center text-base font-normal"
                  disabled
                >
                  <div className="absolute left-4 size-5">
                    {status === "loading" ? (
                      <Loader2 className="h-5 w-5 animate-spin" />
                    ) : (
                      <Check className="h-5 w-5 text-green-600" />
                    )}
                  </div>
                  {status === "loading"
                    ? "Installing..."
                    : "Connected with Slack"}
                </Button>
              </div>
            </motion.div>
          </AnimatePresence>
        </div>
      </div>

      <div className="hidden md:block md:w-1/3 lg:w-1/2 bg-white relative">
        <AnimatePresence mode="wait">
          <motion.div
            className="absolute inset-0 bg-gray-50"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
          />
          <motion.img
            src="/onboarding/four.png"
            alt="Installation Preview"
            className="absolute top-0 left-0 w-auto h-full object-cover object-left pl-8 pt-6 pb-6"
            style={{
              userSelect: "none",
              WebkitUserSelect: "none",
              MozUserSelect: "none",
              msUserSelect: "none",
            }}
            draggable={false}
            onDragStart={(e) => e.preventDefault()}
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: -20 }}
            transition={{
              type: "spring",
              stiffness: 400,
              damping: 30,
              mass: 0.8,
            }}
          />
        </AnimatePresence>
      </div>
    </div>
  );
}

export default function CompletePage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <CompletePageContent />
    </Suspense>
  );
}
