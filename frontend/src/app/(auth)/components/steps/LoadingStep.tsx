// src/components/steps/LoadingStep.tsx
"use client";

import { LOADING_MESSAGES, MESSAGE_INTERVAL } from "@/app/_constants/loading";
import { AnimatePresence, motion } from "framer-motion";
import Link from "next/link";
import { useEffect, useState } from "react";

// Default message for initial render to avoid hydration mismatch
const DEFAULT_MESSAGE = "Preparing your experience";

interface LoadingStepProps {
  message?: string;
}

const LoadingStep = ({ message }: LoadingStepProps) => {
  const [isMounted, setIsMounted] = useState(false);
  const [currentMessage, setCurrentMessage] = useState(message ?? DEFAULT_MESSAGE);

  // Mark component as mounted
  useEffect(() => {
    setIsMounted(true);
  }, []);

  // Start message cycling only after mounting
  useEffect(() => {
    if (!isMounted || message) return;

    const getRandomMessage = () => {
      const messages = LOADING_MESSAGES.filter(msg => msg !== currentMessage);
      return messages[Math.floor(Math.random() * messages.length)];
    };

    const messageInterval = setInterval(() => {
      setCurrentMessage(getRandomMessage());
    }, MESSAGE_INTERVAL);

    return () => clearInterval(messageInterval);
  }, [isMounted, message, currentMessage]);

  const year = new Date().getFullYear();

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.5 }}
      className="min-h-screen bg-white text-black flex flex-col"
    >
      <nav>
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
          <Link
            href="/"
            className="text-xl font-bold hover:scale-105 transition-transform"
          >
            byrd
          </Link>
        </div>
      </nav>

      <main className="flex-grow flex flex-col justify-center items-center px-4 text-center">
        <motion.h1
          initial={{ y: -20, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ delay: 0.2, duration: 0.5 }}
          className="text-6xl font-bold text-black mb-6"
        >
          Just a moment
        </motion.h1>

        <AnimatePresence mode="wait">
          <motion.p
            key={currentMessage}
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            exit={{ y: -20, opacity: 0 }}
            transition={{ duration: 0.3 }}
            className="text-lg text-gray-600 max-w-md mb-8"
          >
            {currentMessage}
          </motion.p>
        </AnimatePresence>

        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{
            delay: 0.5,
            type: "spring",
            stiffness: 200,
            damping: 20,
          }}
        >
          <div className="w-16 h-16 border-4 border-black border-t-transparent rounded-full animate-spin" />
        </motion.div>
      </main>

      <motion.footer
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ delay: 0.7, duration: 0.5 }}
        className="py-8 px-6"
      >
        <div className="max-w-7xl mx-auto">
          <div className="flex justify-between items-center text-sm text-gray-600">
            <p>© {year} byrd. All rights reserved.</p>
            <div className="flex space-x-4">
              <Link href="/terms" className="hover:text-black transition-colors">
                Terms
              </Link>
              <Link href="/privacy" className="hover:text-black transition-colors">
                Privacy
              </Link>
            </div>
          </div>
        </div>
      </motion.footer>
    </motion.div>
  );
};

export default LoadingStep;
