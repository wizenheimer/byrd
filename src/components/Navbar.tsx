"use client";

import { Button } from "@/components/ui/button";
import {
	NavigationMenu,
	NavigationMenuItem,
	NavigationMenuList,
	NavigationMenuTrigger,
} from "@/components/ui/navigation-menu";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { ChevronDown, Menu, X } from "lucide-react";
import React, { useState } from "react";

const navigationItems = {
	product: {
		title: "Product",
		sections: {
			productIntelligence: {
				title: "Product Intelligence",
				showOnMobile: true,
				items: [
					{ name: "Product Launches", href: "/" },
					{ name: "Roadmap Changes", href: "/" },
					{ name: "Feature Releases", href: "/" },
					{ name: "Integration Highlights", href: "/" },
				],
			},
			mediaIntelligence: {
				title: "Media Intelligence",
				showOnMobile: false,
				items: [
					{ name: "Press Release Tracking", href: "/" },
					{ name: "Funding Rounds", href: "/" },
					{ name: "Acquisitions and Mergers", href: "/" },
					{ name: "Leadership Changes", href: "/" },
				],
			},
			competitiveIntelligence: {
				title: "Competitive Intelligence",
				showOnMobile: true,
				items: [
					{ name: "Price Monitoring", href: "/" },
					{ name: "Partnership Briefings", href: "/" },
					{ name: "Positioning Changes", href: "/" },
					{ name: "Promotional Offers", href: "/" },
				],
			},
			customerIntelligence: {
				title: "Customer Intelligence",
				showOnMobile: false,
				items: [
					{ name: "Sentiment Overview", href: "/" },
					{ name: "Review Highlights", href: "/" },
					{ name: "Testimonial Changes", href: "/" },
				],
			},
			socialIntelligence: {
				title: "Social Intelligence",
				showOnMobile: false,
				items: [
					{ name: "Engagement Metrics", href: "/" },
					{ name: "Content Analysis", href: "/" },
				],
			},
			marketingIntelligence: {
				title: "Marketing Intelligence",
				showOnMobile: false,
				items: [
					{ name: "Content Strategy Shifts", href: "/" },
					{ name: "Newsletter Insights", href: "/" },
				],
			},
		},
	},
	resources: {
		title: "Resources",
		sections: {
			resources: {
				title: "Resources",
				showOnMobile: false,
				items: [
					{ name: "Swipe Files", href: "/" },
					{ name: "Ad Library", href: "/" },
					{ name: "Newsletter Library", href: "/" },
					{ name: "Interface Library", href: "/" },
				],
			},
			byrd: {
				title: "Byrd",
				showOnMobile: true,
				items: [
					{ name: "Release Notes", href: "/" },
					{ name: "Blog", href: "/" },
					{ name: "About Us", href: "/" },
				],
			},
			help: {
				title: "Help",
				showOnMobile: true,
				items: [
					{ name: "Slack Community", href: "/" },
					{ name: "Support", href: "/" },
					{ name: "Hire an expert", href: "/" },
					{ name: "System Status", href: "/" },
				],
			},
			integrations: {
				title: "Integrations",
				showOnMobile: false,
				items: [
					{ name: "Slack", href: "/" },
					{ name: "Notion", href: "/" },
					{ name: "Google Workspace", href: "/" },
				],
			},
		},
	},
};

const navAccentStyle =
	"text-gray-600 hover:text-black block w-full rounded-md p-2 hover:bg-accent hover:text-accent-foreground transition-colors";

const ProductDropdown = () => {
	return (
		<div className="absolute top-full left-0 w-full bg-white py-8 px-16 shadow-lg">
			<div className="grid grid-cols-3 gap-8 max-w-7xl mx-auto">
				{Object.entries(navigationItems.product.sections).map(
					([key, section]) => (
						<div key={key}>
							<h3 className="font-semibold mb-4">{section.title}</h3>
							<ul className="space-y-3">
								{section.items.map((item) => (
									<li key={item.name}>
										<a href={item.href} className={navAccentStyle}>
											{item.name}
										</a>
									</li>
								))}
							</ul>
						</div>
					),
				)}
			</div>
		</div>
	);
};

const ResourcesDropdown = () => {
	return (
		<div className="absolute top-full left-0 w-full bg-white py-8 px-16 shadow-lg">
			<div className="grid grid-cols-12 gap-8 max-w-7xl mx-auto">
				<div className="col-span-7 grid grid-cols-2 gap-8">
					{Object.entries(navigationItems.resources.sections).map(
						([key, section]) => (
							<div key={key}>
								<h3 className="font-semibold mb-4">{section.title}</h3>
								<ul className="space-y-3">
									{section.items.map((item) => (
										<li key={item.name}>
											<a href={item.href} className={navAccentStyle}>
												{item.name}
											</a>
										</li>
									))}
								</ul>
							</div>
						),
					)}
				</div>
				<div className="col-span-5">
					<a href="/blog/uber-case-study" className="block group">
						<div className="rounded-xl overflow-hidden">
							<img
								src="/assets/blog-cover.png"
								alt="Blog cover"
								className="w-full h-48 object-cover rounded-xl"
							/>
						</div>
						<h3 className="mt-4 text-lg font-medium group-hover:text-gray-900">
							How Uber established dominance in ride sharing
						</h3>
						<p className="mt-1 text-sm text-gray-600">by Massimo Ruggero</p>
					</a>
				</div>
			</div>
		</div>
	);
};

const MobileMenuItem = ({
	title,
	sections,
}: {
	title: string;
	sections:
		| typeof navigationItems.product.sections
		| typeof navigationItems.resources.sections;
}) => {
	return (
		<Collapsible className="border-b border-gray-200 py-4">
			<CollapsibleTrigger className="flex w-full items-center justify-between">
				<span className="text-lg font-medium">{title}</span>
				<ChevronDown className="h-5 w-5 text-gray-500" />
			</CollapsibleTrigger>
			<CollapsibleContent className="mt-4 space-y-4">
				{Object.entries(sections).map(
					([key, section]) =>
						section.showOnMobile && (
							<div key={key}>
								<h4 className="mb-2">
									<span className="inline-flex items-center justify-center px-4 py-2 mb-4 text-md font-medium bg-gray-100 rounded-full text-gray-800">
										{section.title}
									</span>
								</h4>
								<ul className="space-y-4 pl-4">
									{section.items.map((item) => (
										<li key={item.name}>
											<a
												href={item.href}
												className="block py-1 text-gray-600 hover:text-gray-900"
											>
												{item.name}
											</a>
										</li>
									))}
								</ul>
							</div>
						),
				)}
			</CollapsibleContent>
		</Collapsible>
	);
};

const MobileMenu = () => {
	return (
		<div className="py-4">
			<div className="space-y-4">
				{Object.entries(navigationItems).map(([key, item]) => (
					<MobileMenuItem
						key={key}
						title={item.title}
						sections={item.sections}
					/>
				))}
			</div>
		</div>
	);
};

const Navbar = () => {
	const [activeDropdown, setActiveDropdown] = useState<string | null>(null);
	const [isOverlayVisible, setIsOverlayVisible] = useState(false);
	const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

	const handleDropdownChange = (dropdownName: string | null) => {
		setActiveDropdown(dropdownName);
		setIsOverlayVisible(!!dropdownName);
	};

	return (
		<>
			{isOverlayVisible && (
				<div
					className="fixed inset-0 z-40 bg-black/20 backdrop-blur-sm"
					onClick={() => handleDropdownChange(null)}
					onKeyUp={(e) => {
						if (e.key === "Escape") handleDropdownChange(null);
					}}
					tabIndex={0}
					role="button"
					aria-label="Close overlay"
				/>
			)}

			<div className="relative z-50 w-full bg-background">
				<div className="relative">
					<nav
						className={`${activeDropdown ? "bg-white shadow-sm" : ""} transition-colors duration-200`}
					>
						<div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
							<a href="/" className="text-xl font-bold">
								byrd
							</a>
							<div className="hidden items-center gap-8 md:flex">
								<NavigationMenu>
									<NavigationMenuList className="flex justify-center">
										{Object.entries(navigationItems).map(([key, item]) => (
											<NavigationMenuItem key={key}>
												<NavigationMenuTrigger
													onClick={() =>
														handleDropdownChange(
															activeDropdown === key ? null : key,
														)
													}
													className={`text-base ${activeDropdown === key ? "bg-white" : ""}`}
												>
													{item.title}
												</NavigationMenuTrigger>
											</NavigationMenuItem>
										))}
									</NavigationMenuList>
								</NavigationMenu>
							</div>

							<div className="flex items-center gap-4">
								<Button className="hidden bg-black text-white hover:bg-black/90 md:block">
									Get Started
								</Button>
								<Sheet
									open={isMobileMenuOpen}
									onOpenChange={setIsMobileMenuOpen}
								>
									<SheetTrigger asChild>
										<Button variant="outline" size="icon" className="md:hidden">
											<Menu className="h-6 w-6" />
											<span className="sr-only">Toggle menu</span>
										</Button>
									</SheetTrigger>
									<SheetContent
										side="right"
										className="w-full sm:w-[400px] overflow-y-auto"
									>
										<div className="flex items-center justify-between mb-8">
											<a href="/" className="text-xl font-bold">
												byrd
											</a>
											<Button
												variant="ghost"
												size="icon"
												onClick={() => setIsMobileMenuOpen(false)}
											>
												{/* <X className="h-6 w-6" /> */}
												{/* <span className="sr-only">Close menu</span> */}
											</Button>
										</div>
										<MobileMenu />
										<div className="mt-8">
											<Button className="w-full bg-black text-white hover:bg-black/90">
												Get Started
											</Button>
										</div>
									</SheetContent>
								</Sheet>
							</div>
						</div>
					</nav>

					{activeDropdown === "product" && <ProductDropdown />}
					{activeDropdown === "resources" && <ResourcesDropdown />}
				</div>
			</div>
		</>
	);
};

export default Navbar;
