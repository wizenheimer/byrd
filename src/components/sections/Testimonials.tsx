import SectionHeader from "./block/SectionHeaderBlock";

const TestimonialsSection = () => {
	const headerContent = {
		title: "You're in a good company",
		subtitle:
			"From challenger brands to category kings. From underdogs to top dogs.",
	};
	return (
		<div className="w-full bg-background relative overflow-hidden mt-32 md:mt-48 lg:mt-60">
			{/* Hero Text Content */}
			<SectionHeader {...headerContent} />

			{/* Testimonials Grid */}
			<div className="max-w-6xl mx-auto px-4 grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6 mt-16 md:mt-24 lg:mt-30">
				{/* CTO Testimonial */}
				<div className="bg-gray-100 hover:bg-black rounded-2xl p-6 md:p-8 lg:p-10 flex flex-col justify-between transition-colors duration-500 group">
					<div>
						<p className="text-2xl font-semibold mb-6 leading-tight text-gray-900 group-hover:text-white">
							After trying and ditching countless other tools, I was pretty
							jaded. A month later, I&apos;m eating my words. Have never seen
							our sales and product team so pumped about a tool. Absolute
							must-have.
						</p>
					</div>
					<div>
						<div className="font-semibold text-lg mb-1 text-gray-900 group-hover:text-white">
							CTO
						</div>
						<div className="text-gray-600 group-hover:text-gray-400 text-sm">
							Series D
						</div>
						<div className="text-gray-600 group-hover:text-gray-400 text-sm">
							$81M ARR
						</div>
					</div>
				</div>

				{/* Right Column Testimonials */}
				<div className="space-y-4 md:space-y-6">
					{/* VP Product Testimonial */}
					<div className="bg-gray-100 hover:bg-black rounded-2xl p-6 md:p-8 lg:p-10 transition-colors duration-500 group">
						<div className="mb-6">
							<p className="text-xl font-semibold mb-6 leading-tight text-gray-900 group-hover:text-white">
								I&apos;ve used &apos;em all, and this one takes the cake.
								Doesn&apos;t just add value, multiplies it.
							</p>
						</div>
						<div>
							<div className="font-semibold text-lg mb-1 text-gray-900 group-hover:text-white">
								VP, Product
							</div>
							<div className="text-gray-600 group-hover:text-gray-400 text-sm">
								Series C
							</div>
							<div className="text-gray-600 group-hover:text-gray-400 text-sm">
								$54M ARR
							</div>
						</div>
					</div>

					{/* Founder Testimonial */}
					<div className="bg-gray-100 hover:bg-black rounded-2xl p-6 md:p-8 lg:p-10 transition-colors duration-500 group">
						<div className="mb-6">
							<p className="text-xl font-semibold mb-6 leading-tight text-gray-900 group-hover:text-white">
								CFO nearly had a heart attack when I suggested another tool.
								Now? She&apos;s the one championing it.
							</p>
						</div>
						<div>
							<div className="font-semibold text-lg mb-1 text-gray-900 group-hover:text-white">
								Founder
							</div>
							<div className="text-gray-600 group-hover:text-gray-400 text-sm">
								Bootstrapped
							</div>
							<div className="text-gray-600 group-hover:text-gray-400 text-sm">
								$17M ARR
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	);
};

export default TestimonialsSection;
