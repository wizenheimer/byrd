import Link from "next/link";
import React from "react";

type ContentColumn = {
	title: {
		desktop: string;
		mobile: string;
	};
	description: string;
	linkText: string;
	linkHref: string;
};

type ContentGridProps = {
	mainTitle: {
		desktop: string;
		mobile: string;
	};
	columns: ContentColumn[];
};

const MultiLineText = ({ text }: { text: string }) => (
	<>
		{text.split("\n").map((line, i, arr) => {
			const key = `${line}-${i}`;
			return (
				<React.Fragment key={key}>
					{line}
					{i < arr.length - 1 && <br />}
				</React.Fragment>
			);
		})}
	</>
);

const ContentGrid = ({ mainTitle, columns }: ContentGridProps) => {
	return (
		<div className="max-w-6xl mx-auto px-4 relative mt-12 md:mt-32 lg:mt-48">
			<div className="grid grid-cols-1 md:grid-cols-4 gap-8 md:gap-12">
				{/* Title Column - Hidden on mobile */}
				<div className="hidden md:block">
					<h2 className="text-2xl font-bold leading-tight">
						<span className="hidden md:inline">
							<MultiLineText text={mainTitle.desktop} />
						</span>
						<span className="md:hidden">{mainTitle.mobile}</span>
					</h2>
				</div>

				{/* Content Columns */}
				<div className="col-span-full md:col-span-3">
					<div className="grid grid-cols-1 md:grid-cols-3 gap-12 md:gap-12">
						{columns.map((column) => (
							<div
								key={column.linkHref}
								className="space-y-6 md:grid md:grid-rows-[auto_1fr_auto] md:gap-8 max-w-sm mx-auto md:mx-0"
							>
								<h2 className="text-2xl font-bold">
									<span className="hidden md:inline">
										<MultiLineText text={column.title.desktop} />
									</span>
									<span className="md:hidden">{column.title.mobile}</span>
								</h2>
								<p className="text-gray-600 leading-relaxed">
									{column.description}
								</p>
								<div>
									<Link
										href={column.linkHref}
										className="text-gray-500 font-medium underline underline-offset-4 hover:text-gray-700"
									>
										{column.linkText}
									</Link>
								</div>
							</div>
						))}
					</div>
				</div>
			</div>
		</div>
	);
};

export default ContentGrid;
