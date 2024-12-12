// src/app/(onboarding)/page.tsx
import Footer from "@/components/Footer";
import Navbar from "@/components/Navbar";
import HeroSection from "./components/sections/HeroSection";
import ContextSection from "./components/sections/ContextSection";
import MonitoringSection from "./components/sections/MonitoringSection";
import InboxSection from "./components/sections/InboxSection";
import MediaIntelligenceSection from "./components/sections/MediaIntelligenceSection";
import CustomerIntelligenceSection from "./components/sections/CustomerIntelligenceSection";
import ScaleupSection from "./components/sections/ScaleupSection";
import IntegrationSection from "./components/sections/IntegrationSection";
import TestimonialsSection from "./components/sections/TestimonialSection";
import ClosingNoteSection from "./components/sections/ClosingNoteSection";

export default function Home() {
  return (
    <>
      <Navbar />
      <HeroSection />
      <ContextSection />
      <MonitoringSection />
      <InboxSection />
      <MediaIntelligenceSection />
      <CustomerIntelligenceSection />
      <ScaleupSection />
      <IntegrationSection />
      <TestimonialsSection />
      <ClosingNoteSection />
      <Footer />
    </>
  );
}
