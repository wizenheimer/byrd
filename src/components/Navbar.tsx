"use client";

import React, { useState } from 'react';
import { Button } from "@/components/ui/button";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuList,
  NavigationMenuTrigger,
} from "@/components/ui/navigation-menu";

const navigationItems = {
  productIntelligence: {
    title: "Product Intelligence",
    items: [
      "Product Launches",
      "Roadmap Changes",
      "Feature Releases",
      "Integration Highlights"
    ]
  },
  mediaIntelligence: {
    title: "Media Intelligence",
    items: [
      "Press Release Tracking",
      "Funding Rounds",
      "Acquisitions and Mergers",
      "Leadership Changes"
    ]
  },
  competitiveIntelligence: {
    title: "Competitive Intelligence",
    items: [
      "Price Monitoring",
      "Partnership Briefings",
      "Positioning Changes",
      "Promotional Offers"
    ]
  },
  customerIntelligence: {
    title: "Customer Intelligence",
    items: [
      "Sentiment Overview",
      "Review Highlights",
      "Testimonial Changes"
    ]
  },
  integrations: {
    title: "Integrations",
    items: [
      "Slack",
      "Notion",
      "Google Workspace",
    ]
  },
  socialIntelligence: {
    title: "Social Intelligence",
    items: [
      "Engagement Metrics",
      "Content Analysis"
    ]
  },
  marketingIntelligence: {
    title: "Marketing Intelligence",
    items: [
      "Content Strategy Shifts",
      "Newsletter Insights"
    ]
  }
};

const navAccentStyle = "text-gray-600 hover:text-black block w-full rounded-md p-2 hover:bg-accent hover:text-accent-foreground transition-colors"

const ProductDropdown = () => {
  return (
    <div className="absolute top-full left-0 w-full bg-white py-8 px-16">
      <div className="grid grid-cols-4 gap-8 max-w-7xl mx-auto">
        {/* First Row */}
        <div>
          <h3 className="font-semibold mb-4">{navigationItems.productIntelligence.title}</h3>
          <ul className="space-y-3">
            {navigationItems.productIntelligence.items.map((item) => (
              <li key={item}>
                <a href="/" className={navAccentStyle}>{item}</a>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <h3 className="font-semibold mb-4">{navigationItems.mediaIntelligence.title}</h3>
          <ul className="space-y-3">
            {navigationItems.mediaIntelligence.items.map((item) => (
              <li key={item}>
                <a href="/" className={navAccentStyle}>{item}</a>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <h3 className="font-semibold mb-4">{navigationItems.competitiveIntelligence.title}</h3>
          <ul className="space-y-3">
            {navigationItems.competitiveIntelligence.items.map((item) => (
              <li key={item}>
                <a href="/" className={navAccentStyle}>{item}</a>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <h3 className="font-semibold mb-4">{navigationItems.customerIntelligence.title}</h3>
          <ul className="space-y-3">
            {navigationItems.customerIntelligence.items.map((item) => (
              <li key={item}>
                <a href="/" className={navAccentStyle}>{item}</a>
              </li>
            ))}
          </ul>
        </div>

        {/* Second Row */}
        <div className="mt-8">
          <h3 className="font-semibold mb-4">{navigationItems.integrations.title}</h3>
          <ul className="space-y-3">
            {navigationItems.integrations.items.map((item) => (
              <li key={item}>
                <a href="/" className={navAccentStyle}>{item}</a>
              </li>
            ))}
          </ul>
        </div>
        <div className="mt-8">
          <h3 className="font-semibold mb-4">{navigationItems.socialIntelligence.title}</h3>
          <ul className="space-y-3">
            {navigationItems.socialIntelligence.items.map((item) => (
              <li key={item}>
                <a href="/" className={navAccentStyle}>{item}</a>
              </li>
            ))}
          </ul>
        </div>
        <div className="mt-8">
          <h3 className="font-semibold mb-4">{navigationItems.marketingIntelligence.title}</h3>
          <ul className="space-y-3">
            {navigationItems.marketingIntelligence.items.map((item) => (
              <li key={item}>
                <a href="/" className={navAccentStyle}>{item}</a>
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
    <div className="absolute top-full left-0 w-full bg-white py-8 px-16">
      <div className="grid grid-cols-12 gap-8 max-w-7xl mx-auto">
        {/* Left Section - Navigation Links */}
        <div className="col-span-7 grid grid-cols-2 gap-8">
          <div>
            <h3 className="font-semibold mb-4">Resources</h3>
            <ul className="space-y-3">
              <li><a href="/" className={navAccentStyle}>Swipe Files</a></li>
              <li><a href="/" className={navAccentStyle}>Ad Library</a></li>
              <li><a href="/" className={navAccentStyle}>Newsletter Library</a></li>
              <li><a href="/" className={navAccentStyle}>Interface Library</a></li>
            </ul>

            <h3 className="font-semibold mb-4 mt-8">Byrd</h3>
            <ul className="space-y-3">
              <li><a href="/" className={navAccentStyle}>Release Notes</a></li>
              <li><a href="/" className={navAccentStyle}>Blog</a></li>
              <li><a href="/" className={navAccentStyle}>About Us</a></li>
            </ul>
          </div>

          <div>
            <h3 className="font-semibold mb-4">Help</h3>
            <ul className="space-y-3">
              <li><a href="/" className={navAccentStyle}>Slack Community</a></li>
              <li><a href="/" className={navAccentStyle}>Support</a></li>
              <li><a href="/" className={navAccentStyle}>Hire an expert</a></li>
              <li><a href="/" className={navAccentStyle}>System Status</a></li>
            </ul>
          </div>
        </div>

        {/* Right Section - Featured Content */}
        <div className="col-span-5">
          <a href="/" className="block group">
            <div className="rounded-xl overflow-hidden">
              <img 
                src="assets/blog-cover.png" 
                alt="Ocean waves" 
                className="w-full h-48 object-cover rounded-xl"
              />
            </div>
            <h3 className="mt-4 text-lg font-medium group-hover:text-gray-900">
              How Uber established dominance in ride sharing
            </h3>
            <p className="mt-1 text-sm text-gray-600">
              by Massimo Ruggero
            </p>
          </a>
        </div>
      </div>
    </div>
  );
};

const Navbar = () => {
  const [activeDropdown, setActiveDropdown] = useState<string | null>(null);

  return (
    <>
    <div className="w-full bg-background">
    <div className="relative">
      <nav className={`${activeDropdown ? 'border-none bg-white' : ''}`}>
      <div className="max-w-7xl mx-auto px-4 flex items-center justify-between h-16">
        <a href="/" className="font-bold text-xl">byrd</a>
        <div className="flex items-center gap-8">
        <NavigationMenu>
          <NavigationMenuList className="flex justify-center">
          <NavigationMenuItem>
            <NavigationMenuTrigger 
            onClick={() => setActiveDropdown(activeDropdown === 'product' ? null : 'product')}
            className="text-base"
            >
            Product
            </NavigationMenuTrigger>
          </NavigationMenuItem>
          <NavigationMenuItem>
            <NavigationMenuTrigger 
            onClick={() => setActiveDropdown(activeDropdown === 'resources' ? null : 'resources')}
            className="text-base"
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
      
      {activeDropdown === 'product' && <ProductDropdown />}
      {activeDropdown === 'resources' && <ResourcesDropdown />}
    </div>
    </div>
    </>
  );
};

export default Navbar;