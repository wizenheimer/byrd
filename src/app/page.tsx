import BigPictureSection from "@/components/BigPictureSection";
import CompetitionSection from "@/components/Competition";
import DifferentiationSection from "@/components/DifferentiationSection";
import ExitCTA from "@/components/ExitCTA";
import HeroSection from "@/components/HeroSection";
import LastOne from "@/components/LastOne";
import MonitoringSection from "@/components/MonitoringSection";
import TestimonialsSection from "@/components/Testimonials";

export default function Home() {
	return (
		<>
			<HeroSection />
			<MonitoringSection />
			<BigPictureSection />
			<LastOne />
			<CompetitionSection />
			<DifferentiationSection />
			<TestimonialsSection />
			<ExitCTA />
		</>
	);
}
