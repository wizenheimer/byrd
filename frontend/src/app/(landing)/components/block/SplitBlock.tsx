import AnimatedImage from "@/components/AnimatedImage";
import Link from "next/link";

type ColumnContent = {
	imageSrc: string;
	title: string;
	paragraphs: string[];
	linkText: string;
	linkHref: string;
};

type SplitBlockProps = {
	leftColumn: ColumnContent;
	rightColumn: ColumnContent;
};

const SplitBlock = ({ leftColumn, rightColumn }: SplitBlockProps) => {
	const renderColumn = (content: ColumnContent) => (
		<div>
			<div className="relative mb-16">
				<AnimatedImage imageSrc={content.imageSrc} />
			</div>
			<div className="space-y-6">
				<h2 className="text-3xl font-bold">{content.title}</h2>
				{content.paragraphs.map((paragraph, index) => (
					// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
					<p key={index} className="text-gray-600">
						{paragraph}
					</p>
				))}
			</div>
			<div className="mt-12">
				<Link
					href={content.linkHref}
					className="text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
				>
					{content.linkText}
				</Link>
			</div>
		</div>
	);

	return (
		<div className="max-w-7xl mx-auto">
			<div className="grid grid-cols-1 lg:grid-cols-2 gap-24 px-8 py-12">
				{/* Left Column */}
				{renderColumn(leftColumn)}

				{/* Right Column */}
				{renderColumn(rightColumn)}
			</div>
		</div>
	);
};

export default SplitBlock;
