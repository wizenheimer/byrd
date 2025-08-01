"use client";

import { Button } from "@/components/ui/button";
import { Github } from "lucide-react";
import Link from "next/link";

const Navbar = () => {
  return (
    <div className="relative z-50 w-full bg-background">
      <nav>
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
          <a href="/" className="text-xl font-bold">
            byrd
          </a>
        </div>
      </nav>
    </div>
  );
};

export default Navbar;
