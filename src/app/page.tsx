import Footer from "@/components/Footer";
import Navbar from "@/components/Navbar";
import ClosingNoteSection from "@/components/sections/ClosingNoteSection";
import ContextSection from "@/components/sections/ContextSection";
import CustomerIntelligenceSection from "@/components/sections/CustomerIntelligenceSection";
import HeroSection from "@/components/sections/HeroSection";
import InboxSection from "@/components/sections/InboxSection";
import IntegrationSection from "@/components/sections/IntegrationSection";
import MediaIntelligenceSection from "@/components/sections/MediaIntelligenceSection";
import MonitoringSection from "@/components/sections/MonitoringSection";
import ScaleupSection from "@/components/sections/ScaleupSection";
import TestimonialsSection from "@/components/sections/TestimonialSection";

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
