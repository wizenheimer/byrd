// src/components/block/FeatureGridBlock.tsx
import Link from "next/link";
import React from "react";

type GridColumn = {
	title: {
		desktop: string;
		mobile: string;
	};
	description: string;
	link: {
		text: {
			desktop: string;
			mobile: string;
		};
		href: string;
	};
};

type FeatureGridProps = {
	columns: GridColumn[];
	className?: string;
};

const MultiLineText = ({ text }: { text: string }) => (
	<>
		{text.split("\n").map((line, i, arr) => (
			// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
			<React.Fragment key={i}>
				{line}
				{i < arr.length - 1 && <br />}
			</React.Fragment>
		))}
	</>
);

const FeatureGrid = ({ columns, className = "" }: FeatureGridProps) => {
	return (
		<div
			className={`max-w-6xl mx-auto px-4 relative mt-12 md:mt-32 lg:mt-48 ${className}`}
		>
			<div className="grid grid-cols-1 md:grid-cols-3 gap-12 md:gap-24">
				{columns.map((column, index) => (
					<div
						// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
						key={index}
						className="space-y-6 md:grid md:grid-rows-[auto_1fr_auto] md:gap-8 max-w-sm mx-auto md:mx-0"
					>
						<h2 className="text-2xl font-bold text-center md:text-left">
							<span className="hidden md:inline">
								<MultiLineText text={column.title.desktop} />
							</span>
							<span className="md:hidden">{column.title.mobile}</span>
						</h2>
						<p className="text-gray-600 leading-relaxed">
							{column.description}
						</p>
						<div className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700">
							<Link href={column.link.href} className="hidden md:inline">
								<MultiLineText text={column.link.text.desktop} />
							</Link>
							<Link href={column.link.href} className="md:hidden">
								{column.link.text.mobile}
							</Link>
						</div>
					</div>
				))}
			</div>
		</div>
	);
};

export default FeatureGrid;
