type Props = {
	imageSrc: string;
	altText?: string;
};

const OnboardingPreviewPane = ({ imageSrc, altText }: Props) => {
	return (
		<div className="hidden md:block md:w-1/3 lg:w-1/2 bg-gray-50 overflow-hidden relative">
			<img
				src={imageSrc}
				alt={altText || "Dashboard Preview"}
				className="absolute top-0 left-0 w-auto h-full object-cover object-left pl-8 pt-6 pb-6"
				style={{
					userSelect: "none",
					WebkitUserSelect: "none",
					MozUserSelect: "none",
					msUserSelect: "none",
				}}
				draggable={false}
				onDragStart={(e) => e.preventDefault()}
			/>
		</div>
	);
};

export default OnboardingPreviewPane;
