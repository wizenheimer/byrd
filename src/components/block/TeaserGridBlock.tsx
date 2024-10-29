import Link from "next/link";
import React from "react";

type ContentColumn = {
	title: {
		desktop: string;
		mobile?: string; // Optional since first column doesn't have mobile version
	};
	description?: string; // Optional since first column doesn't have description
	linkText?: string; // Optional since first column doesn't have link
	linkHref?: string; // Optional since first column doesn't have link
};

type ThreeColumnSectionProps = {
	columns: [ContentColumn, ContentColumn, ContentColumn]; // Exactly 3 columns
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

const ThreeColumnSection = ({ columns }: ThreeColumnSectionProps) => {
	return (
		<div className="max-w-6xl mx-auto px-4 relative mt-12 md:mt-32 lg:mt-48">
			<div className="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-20">
				{/* Column 1 */}
				<div className="md:block grid grid-rows-[1fr] items-start">
					<h2 className="text-3xl font-semibold">
						<span className="hidden md:inline">
							<MultiLineText text={columns[0].title.desktop} />
						</span>
					</h2>
				</div>

				{/* Columns 2 and 3 */}
				{columns.slice(1).map((column, index) => (
					<div
						// biome-ignore lint/suspicious/noArrayIndexKey: <explanation>
						key={index}
						className="grid grid-rows-[auto_1fr_auto] gap-6 md:gap-8 max-w-sm mx-auto md:mx-0"
					>
						<h2 className="text-2xl font-semibold">
							<span className="hidden md:inline">
								<MultiLineText text={column.title.desktop} />
							</span>
							<span className="md:hidden">{column.title.mobile}</span>
						</h2>
						<p className="text-base font-medium text-gray-600 leading-relaxed">
							{column.description}
						</p>
						<div>
							<Link
								href={column.linkHref || "/"}
								className="inline-block text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
							>
								{column.linkText}
							</Link>
						</div>
					</div>
				))}
			</div>
		</div>
	);
};

export default ThreeColumnSection;
