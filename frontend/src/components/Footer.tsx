import { Github, Inbox, Linkedin } from "lucide-react";
import Link from "next/link";

const footerData = {
	legalLinks: [
		{ label: "Privacy Notice", href: "/privacy" },
		{ label: "Terms of Service", href: "/terms" },
		{ label: "End User Agreement", href: "/eula" },
	],
};

const Footer = () => {
	return (
		<footer className="w-full bg-gray-100 py-8 md:py-16 mt-20">
			<div className="max-w-7xl mx-auto px-4">
				{/* Bottom Section */}
				<div className="flex flex-col sm:flex-row justify-between gap-8">
					{/* Left side with social icons and copyright */}
					<div className="order-2 sm:order-1">
						{/* Social Links */}
						<div className="flex gap-6 mb-4 justify-center sm:justify-start">
							<Link
								href="https://www.github.com/byrdlabs"
								className="text-gray-600 hover:text-gray-900 transition-colors"
							>
								<Github className="w-5 h-5" />
							</Link>
							<Link
								href="https://www.linkedin.com/company/byrdhq"
								className="text-gray-600 hover:text-gray-900 transition-colors"
							>
								<Linkedin className="w-5 h-5" />
							</Link>
							<Link
								href="mailto:hey@byrdhq.com"
								className="text-gray-600 hover:text-gray-900 transition-colors"
							>
								<Inbox className="w-5 h-5" />
							</Link>
						</div>
						{/* Copyright */}
						<div className="text-gray-600 text-center sm:text-left">
							© 2024 ByrdLabs
						</div>
					</div>

					{/* Right side with stacked legal links */}
					<div className="flex flex-col sm:flex-row gap-4 sm:gap-8 order-1 sm:order-2 items-center sm:items-start">
						{footerData.legalLinks.map((link) => (
							<Link
								key={link.href}
								href={link.href}
								className="text-gray-600 hover:text-gray-900 transition-colors"
							>
								{link.label}
							</Link>
						))}
					</div>
				</div>
			</div>
		</footer>
	);
};

export default Footer;
