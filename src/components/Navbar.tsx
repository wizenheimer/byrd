"use client";

import { Button } from "@/components/ui/button";
import {
	NavigationMenu,
	NavigationMenuItem,
	NavigationMenuList,
	NavigationMenuTrigger,
} from "@/components/ui/navigation-menu";
import React, { useState } from "react";
import Image from "next/image";

const navigationItems = {
	productIntelligence: {
		title: "Product Intelligence",
		items: [
			{ name: "Product Launches", href: "/product/launches" },
			{ name: "Roadmap Changes", href: "/product/roadmap" },
			{ name: "Feature Releases", href: "/product/features" },
			{ name: "Integration Highlights", href: "/product/integrations" },
		],
	},
	mediaIntelligence: {
		title: "Media Intelligence",
		items: [
			{ name: "Press Release Tracking", href: "/media/press" },
			{ name: "Funding Rounds", href: "/media/funding" },
			{ name: "Acquisitions and Mergers", href: "/media/acquisitions" },
			{ name: "Leadership Changes", href: "/media/leadership" },
		],
	},
	competitiveIntelligence: {
		title: "Competitive Intelligence",
		items: [
			{ name: "Price Monitoring", href: "/competitive/pricing" },
			{ name: "Partnership Briefings", href: "/competitive/partnerships" },
			{ name: "Positioning Changes", href: "/competitive/positioning" },
			{ name: "Promotional Offers", href: "/competitive/promotions" },
		],
	},
	customerIntelligence: {
		title: "Customer Intelligence",
		items: [
			{ name: "Sentiment Overview", href: "/customer/sentiment" },
			{ name: "Review Highlights", href: "/customer/reviews" },
			{ name: "Testimonial Changes", href: "/customer/testimonials" },
		],
	},
	integrations: {
		title: "Integrations",
		items: [
			{ name: "Slack", href: "/integrations/slack" },
			{ name: "Notion", href: "/integrations/notion" },
			{ name: "Google Workspace", href: "/integrations/google" },
		],
	},
	socialIntelligence: {
		title: "Social Intelligence",
		items: [
			{ name: "Engagement Metrics", href: "/social/engagement" },
			{ name: "Content Analysis", href: "/social/content" },
		],
	},
	marketingIntelligence: {
		title: "Marketing Intelligence",
		items: [
			{ name: "Content Strategy Shifts", href: "/marketing/strategy" },
			{ name: "Newsletter Insights", href: "/marketing/newsletter" },
		],
	},
};

const resourcesItems = {
	resources: {
		title: "Resources",
		items: [
			{ name: "Swipe Files", href: "/resources/swipe-files" },
			{ name: "Ad Library", href: "/resources/ad-library" },
			{ name: "Newsletter Library", href: "/resources/newsletter-library" },
			{ name: "Interface Library", href: "/resources/interface-library" },
		],
	},
	byrd: {
		title: "Byrd",
		items: [
			{ name: "Release Notes", href: "/release-notes" },
			{ name: "Blog", href: "/blog" },
			{ name: "About Us", href: "/about" },
		],
	},
	help: {
		title: "Help",
		items: [
			{ name: "Slack Community", href: "/community" },
			{ name: "Support", href: "/support" },
			{ name: "Hire an expert", href: "/experts" },
			{ name: "System Status", href: "/status" },
		],
	},
};

const navAccentStyle =
	"text-gray-600 hover:text-black block w-full rounded-md p-2 hover:bg-accent hover:text-accent-foreground transition-colors";

const ProductDropdown = () => {
	return (
		<div className="absolute top-full left-0 w-full bg-white py-8 px-16 shadow-lg">
			<div className="grid grid-cols-4 gap-8 max-w-7xl mx-auto">
				{/* First Row */}
				<div>
					<h3 className="font-semibold mb-4">
						{navigationItems.productIntelligence.title}
					</h3>
					<ul className="space-y-3">
						{navigationItems.productIntelligence.items.map((item) => (
							<li key={item.name}>
								<a href={item.href} className={navAccentStyle}>
									{item.name}
								</a>
							</li>
						))}
					</ul>
				</div>
				<div>
					<h3 className="font-semibold mb-4">
						{navigationItems.mediaIntelligence.title}
					</h3>
					<ul className="space-y-3">
						{navigationItems.mediaIntelligence.items.map((item) => (
							<li key={item.name}>
								<a href={item.href} className={navAccentStyle}>
									{item.name}
								</a>
							</li>
						))}
					</ul>
				</div>
				<div>
					<h3 className="font-semibold mb-4">
						{navigationItems.competitiveIntelligence.title}
					</h3>
					<ul className="space-y-3">
						{navigationItems.competitiveIntelligence.items.map((item) => (
							<li key={item.name}>
								<a href={item.href} className={navAccentStyle}>
									{item.name}
								</a>
							</li>
						))}
					</ul>
				</div>
				<div>
					<h3 className="font-semibold mb-4">
						{navigationItems.customerIntelligence.title}
					</h3>
					<ul className="space-y-3">
						{navigationItems.customerIntelligence.items.map((item) => (
							<li key={item.name}>
								<a href={item.href} className={navAccentStyle}>
									{item.name}
								</a>
							</li>
						))}
					</ul>
				</div>

				{/* Second Row */}
				<div className="mt-8">
					<h3 className="font-semibold mb-4">
						{navigationItems.integrations.title}
					</h3>
					<ul className="space-y-3">
						{navigationItems.integrations.items.map((item) => (
							<li key={item.name}>
								<a href={item.href} className={navAccentStyle}>
									{item.name}
								</a>
							</li>
						))}
					</ul>
				</div>
				<div className="mt-8">
					<h3 className="font-semibold mb-4">
						{navigationItems.socialIntelligence.title}
					</h3>
					<ul className="space-y-3">
						{navigationItems.socialIntelligence.items.map((item) => (
							<li key={item.name}>
								<a href={item.href} className={navAccentStyle}>
									{item.name}
								</a>
							</li>
						))}
					</ul>
				</div>
				<div className="mt-8">
					<h3 className="font-semibold mb-4">
						{navigationItems.marketingIntelligence.title}
					</h3>
					<ul className="space-y-3">
						{navigationItems.marketingIntelligence.items.map((item) => (
							<li key={item.name}>
								<a href={item.href} className={navAccentStyle}>
									{item.name}
								</a>
							</li>
						))}
					</ul>
				</div>
			</div>
		</div>
	);
};

const ResourcesDropdown = () => {
	return (
		<div className="absolute top-full left-0 w-full bg-white py-8 px-16 shadow-lg">
			<div className="grid grid-cols-12 gap-8 max-w-7xl mx-auto">
				{/* Left Section - Navigation Links */}
				<div className="col-span-7 grid grid-cols-2 gap-8">
					<div>
						<h3 className="font-semibold mb-4">
							{resourcesItems.resources.title}
						</h3>
						<ul className="space-y-3">
							{resourcesItems.resources.items.map((item) => (
								<li key={item.name}>
									<a href={item.href} className={navAccentStyle}>
										{item.name}
									</a>
								</li>
							))}
						</ul>

						<h3 className="font-semibold mb-4 mt-8">
							{resourcesItems.byrd.title}
						</h3>
						<ul className="space-y-3">
							{resourcesItems.byrd.items.map((item) => (
								<li key={item.name}>
									<a href={item.href} className={navAccentStyle}>
										{item.name}
									</a>
								</li>
							))}
						</ul>
					</div>

					<div>
						<h3 className="font-semibold mb-4">{resourcesItems.help.title}</h3>
						<ul className="space-y-3">
							{resourcesItems.help.items.map((item) => (
								<li key={item.name}>
									<a href={item.href} className={navAccentStyle}>
										{item.name}
									</a>
								</li>
							))}
						</ul>
					</div>
				</div>

				{/* Right Section - Featured Content */}
				<div className="col-span-5">
					<a href="/blog/uber-case-study" className="block group">
						<div className="rounded-xl overflow-hidden">
							<Image
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

const Navbar = () => {
	const [activeDropdown, setActiveDropdown] = useState<string | null>(null);

	// Add overlay state
	const [isOverlayVisible, setIsOverlayVisible] = useState(false);

	// Update dropdown handler to manage overlay
	const handleDropdownChange = (dropdownName: string | null) => {
		setActiveDropdown(dropdownName);
		setIsOverlayVisible(!!dropdownName);
	};

	return (
		<>
			{/* Overlay */}
			{isOverlayVisible && (
				<div
					className="fixed inset-0 bg-black/20 z-40 backdrop-blur-sm"
					onClick={() => handleDropdownChange(null)}
					onKeyUp={(e) => {
						if (e.key === "Escape") handleDropdownChange(null);
					}}
					tabIndex={0}
					role="button"
					aria-label="Close overlay"
				/>
			)}

			{/* Main navbar wrapper */}
			<div className="top-0 w-full bg-background relative z-50">
				<div className="relative">
					<nav
						className={`${activeDropdown ? "bg-white shadow-sm" : ""} transition-colors duration-200`}
					>
						<div className="max-w-7xl mx-auto px-4 flex items-center justify-between h-16">
							<a href="/" className="font-bold text-xl">
								byrd
							</a>
							<div className="flex items-center gap-8">
								<NavigationMenu>
									<NavigationMenuList className="flex justify-center">
										<NavigationMenuItem>
											<NavigationMenuTrigger
												onClick={() =>
													handleDropdownChange(
														activeDropdown === "product" ? null : "product",
													)
												}
												className={`text-base ${activeDropdown === "product" ? "bg-white" : ""}`}
											>
												Product
											</NavigationMenuTrigger>
										</NavigationMenuItem>
										<NavigationMenuItem>
											<NavigationMenuTrigger
												onClick={() =>
													handleDropdownChange(
														activeDropdown === "resources" ? null : "resources",
													)
												}
												className={`text-base ${activeDropdown === "resources" ? "bg-white" : ""}`}
											>
												Resources
											</NavigationMenuTrigger>
										</NavigationMenuItem>
									</NavigationMenuList>
								</NavigationMenu>
							</div>

							<Button className="bg-black text-white hover:bg-black/90">
								Get Started
							</Button>
						</div>
					</nav>

					{/* Dropdowns */}
					{activeDropdown === "product" && <ProductDropdown />}
					{activeDropdown === "resources" && <ResourcesDropdown />}
				</div>
			</div>
		</>
	);
};

export default Navbar;
