import { ArrowUpRight } from "lucide-react";
import Link from "next/link";

export default function WaitlistScreen() {
	return (
		<div className="min-h-screen bg-white text-black flex flex-col">
			<nav>
				<div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
					<a href="/" className="text-xl font-bold">
						byrd
					</a>
				</div>
			</nav>

			<main className="flex-grow flex flex-col justify-center items-center px-4 text-center">
				<div className="flex items-center mb-4">
					<h1 className="text-6xl font-bold text-black">Stay tuned</h1>
					<ArrowUpRight className="w-12 h-12 ml-2 text-black" />
				</div>
				<p className="text-lg text-gray-600 max-w-md">
					Rolling out access to 7 new teams every month
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
