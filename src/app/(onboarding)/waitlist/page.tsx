// src/app/(onboarding)/waitlist/page.tsx
"use client";

import LoadingStep from "@/components/steps/LoadingStep";
import { UserButton, useUser } from "@clerk/nextjs";
import { ArrowUpRight } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function WaitlistScreen() {
  const { isLoaded, isSignedIn } = useUser();
  const router = useRouter();

  useEffect(() => {
    if (isLoaded && !isSignedIn) {
      router.push("/get-started");
    }
  }, [isLoaded, isSignedIn, router]);

  if (!isLoaded || !isSignedIn) {
    return <LoadingStep />;
  }

  return (
    <div className="min-h-screen bg-white text-black flex flex-col">
      <nav>
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
          <a href="/" className="text-xl font-bold">
            byrd
          </a>
          <div className="flex items-center space-x-4">
            <UserButton />
          </div>
        </div>
      </nav>

      <main className="flex-grow flex flex-col justify-center items-center px-4 text-center">
        <div className="flex items-center mb-4">
          <h1 className="text-6xl font-bold text-black">Stay tuned</h1>
          <ArrowUpRight className="w-12 h-12 ml-2 text-black" />
        </div>
        <p className="text-md text-gray-600 max-w-md">
          {"You're on the list! We'll reach out soon."}
        </p>
      </main>

      <footer className="py-8 px-6">
        <div className="max-w-7xl mx-auto">
          <div className="flex justify-between items-center text-sm text-gray-600">
            <p>© 2024 byrd. All rights reserved.</p>
            <div className="flex space-x-4">
              <Link href="/terms" className="hover:text-black">
                Terms
              </Link>
              <Link href="/privacy" className="hover:text-black">
                Privacy
              </Link>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
