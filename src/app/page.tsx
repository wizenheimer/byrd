import BigPictureSection from "@/components/sections/BigPictureSection";

import Navbar from "@/components/Navbar";
import HeroSection from "@/components/sections/HeroSection";
import MonitoringSection from "@/components/sections/MonitoringSection";
import LastOne from "@/components/sections/LastOne";
import CompetitionSection from "@/components/sections/Competition";
import DifferentiationSection from "@/components/sections/DifferentiationSection";
import TestimonialsSection from "@/components/sections/Testimonials";
import ExitCTA from "@/components/sections/ExitCTA";
import Footer from "@/components/Footer";

export default function Home() {
	return (
		<>
			<Navbar />
			<HeroSection />
			<MonitoringSection />
			<BigPictureSection />
			<LastOne />
			<CompetitionSection />
			<DifferentiationSection />
			<TestimonialsSection />
			<ExitCTA />
			<Footer />
		</>
	);
}
