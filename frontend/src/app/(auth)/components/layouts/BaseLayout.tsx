// src/components/layouts/BaseLayout.tsx
import Link from "next/link";
import { type ReactNode } from "react";
import { motion } from "framer-motion";

interface BaseLayoutProps {
  children: ReactNode;
  header?: ReactNode;
}

export function BaseLayout({ children, header }: BaseLayoutProps) {
  const currentYear = new Date().getFullYear();

  return (
    <div className="min-h-screen bg-white text-black flex flex-col">
      <nav>
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
          <Link href="/" className="text-xl font-bold hover:scale-105 transition-transform">
            byrd
          </Link>
          {header}
        </div>
      </nav>

      <main className="flex-grow flex flex-col justify-center items-center px-4 text-center">
        {children}
      </main>

      <motion.footer
        initial={{ y: 20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ delay: 0.7, duration: 0.5 }}
        className="py-8 px-6"
      >
        <div className="max-w-7xl mx-auto">
          <div className="flex justify-between items-center text-sm text-gray-600">
            <p>© {currentYear} byrd. All rights reserved.</p>
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
    </div>
  );
}
