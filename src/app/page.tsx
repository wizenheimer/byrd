import Navbar from "@/components/Navbar";
import HeroSection from "@/components/sections/HeroSection";
import ContextSection from "@/components/sections/ContextSection";
import TestimonialsSection from "@/components/sections/TestimonialSection";
import MonitoringSection from "@/components/sections/MonitoringSection";
import InboxSection from "@/components/sections/InboxSection";
import MediaIntelligenceSection from "@/components/sections/MediaIntelligenceSection";
import CustomerIntelligenceSection from "@/components/sections/CustomerIntelligenceSection";
import IntegrationSection from "@/components/sections/IntegrationSection";
import ScaleupSection from "@/components/sections/ScaleupSection";
import ClosingNoteSection from "@/components/sections/ClosingNoteSection";
import Footer from "@/components/Footer";

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
			<IntegrationSection />
			<ScaleupSection />
			<TestimonialsSection />
			<ClosingNoteSection />
			<Footer />
		</>
	);
}
