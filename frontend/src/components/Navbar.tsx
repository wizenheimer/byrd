"use client";

import { Button } from "@/components/ui/button";
import Link from "next/link";

const Navbar = () => {
	return (
		<div className="relative z-50 w-full bg-background">
			<nav>
				<div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
					<a href="/" className="text-xl font-bold">
						byrd
					</a>
					<Link href="/onboarding">
						<Button className="bg-black text-white hover:bg-black/90">
							Get started
						</Button>
					</Link>
				</div>
			</nav>
		</div>
	);
};

export default Navbar;
