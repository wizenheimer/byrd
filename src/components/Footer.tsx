import type React from "react";
import Link from "next/link";
import { Github, Linkedin, Slack } from "lucide-react";
import { Separator } from "@/components/ui/separator";

const footerData = {
	mainSections: [
		{
			title: "Product Intelligence",
			links: [
				{ label: "Product Launches", href: "/" },
				{ label: "Roadmap Changes", href: "/" },
				{ label: "Feature Releases", href: "/" },
				{ label: "Integration Highlights", href: "/" },
				{ label: "Stack Updates", href: "/" },
			],
		},
		{
			title: "Media Intelligence",
			links: [
				{ label: "Press Release Tracking", href: "/" },
				{ label: "Funding Rounds", href: "/" },
				{ label: "Acquisitions and Mergers", href: "/" },
				{ label: "Geographic Expansion", href: "/" },
				{ label: "Leadership Changes", href: "/" },
			],
		},
		{
			title: "Marketing Intelligence",
			links: [
				{ label: "Content Strategy Shifts", href: "/" },
				{ label: "Newsletter Insights", href: "/" },
			],
		},
	],
	rightSection: {
		title: "Byrd",
		links: [
			{ label: "About", href: "/" },
			{ label: "Careers", href: "/" },
			{ label: "Pricing", href: "/" },
		],
	},
	bottomMainSections: [
		{
			title: "Competitive Intelligence",
			links: [
				{ label: "Price Monitoring", href: "/" },
				{ label: "Partnership Briefings", href: "/" },
				{ label: "Positioning Changes", href: "/" },
				{ label: "Promotional Offers", href: "/" },
			],
		},
		{
			title: "Customer Intelligence",
			links: [
				{ label: "Sentiment Overview", href: "/" },
				{ label: "Review Highlights", href: "/" },
				{ label: "Testimonial Changes", href: "/" },
			],
		},
		{
			title: "Social Intelligence",
			links: [
				{ label: "Engagement Metrics", href: "/" },
				{ label: "Content Analysis", href: "/" },
				{ label: "Campaign Insights", href: "/" },
			],
		},
	],
	bottomRightSection: {
		title: "Resources",
		links: [
			{ label: "Community", href: "/" },
			{ label: "Support", href: "/" },
			{ label: "System Status", href: "/" },
		],
	},
	legalLinks: [
		{ label: "Privacy Notice", href: "/" },
		{ label: "Terms of Service", href: "/" },
		{ label: "End User Agreement", href: "/" },
	],
};

interface FooterSectionProps {
	title: string;
	links: { label: string; href: string }[];
}

const FooterSection: React.FC<FooterSectionProps> = ({ title, links }) => (
	<div>
		<h3 className="font-semibold mb-4">{title}</h3>
		<ul className="space-y-3">
			{links.map((link) => (
				<li key={link.href}>
					<Link
						href={link.href}
						className="text-gray-600 hover:text-gray-900 transition-colors"
					>
						{link.label}
					</Link>
				</li>
			))}
		</ul>
	</div>
);

const Footer = () => {
	return (
		<footer className="w-full bg-gray-100 py-16 mt-20">
			<div className="max-w-7xl mx-auto px-4">
				{/* Top Footer Section */}
				<div className="flex justify-between mb-16">
					<div className="grid grid-cols-3 gap-24">
						{footerData.mainSections.map((section) => (
							<FooterSection
								key={section.title}
								title={section.title}
								links={section.links}
							/>
						))}
					</div>
					<div className="min-w-[200px]">
						<FooterSection
							title={footerData.rightSection.title}
							links={footerData.rightSection.links}
						/>
					</div>
				</div>

				{/* Bottom Footer Section */}
				<div className="flex justify-between mb-16">
					<div className="grid grid-cols-3 gap-24">
						{footerData.bottomMainSections.map((section) => (
							<FooterSection
								key={section.title}
								title={section.title}
								links={section.links}
							/>
						))}
					</div>
					<div className="min-w-[200px]">
						<FooterSection
							title={footerData.bottomRightSection.title}
							links={footerData.bottomRightSection.links}
						/>
					</div>
				</div>

				<Separator className="mb-8" />

				{/* Bottom Section */}
				<div className="flex justify-between">
					{/* Left side with social icons and copyright */}
					<div>
						{/* Social Links */}
						<div className="flex gap-6 mb-4">
							<Link
								href="/"
								className="text-gray-600 hover:text-gray-900 transition-colors"
							>
								<Github className="w-5 h-5" />
							</Link>
							<Link
								href="/"
								className="text-gray-600 hover:text-gray-900 transition-colors"
							>
								<Linkedin className="w-5 h-5" />
							</Link>
							<Link
								href="/"
								className="text-gray-600 hover:text-gray-900 transition-colors"
							>
								<Slack className="w-5 h-5" />
							</Link>
						</div>
						{/* Copyright */}
						<div className="text-gray-600">© 2024 ByrdLabs</div>
					</div>

					{/* Right side with stacked legal links */}
					<div className="flex flex-col space-y-4 min-w-[200px]">
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
