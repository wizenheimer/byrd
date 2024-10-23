interface PreviewContainerProps {
	imageSrc: string;
	caption?: string;
}

const PreviewContainer: React.FC<PreviewContainerProps> = ({
	imageSrc,
	caption,
}) => {
	return (
		<>
			<div className="relative bg-gray-100 rounded-md p-2 min-h-[800px] overflow-hidden">
				{/* <div className="relative bg-white rounded-2xl shadow-sm overflow-hidden"> */}
				<div className="w-fit h-fit rounded-lg overflow-hidden border">
					<img
						src={imageSrc}
						alt=""
						className="object-contain absolute -right-[15%] -bottom-[20%] rounded-lg h-[900px] overflow-hidden"
					/>
				</div>
			</div>
			<div className="relative mt-4 text-sm text-gray-600 px-1">{caption}</div>
		</>
	);
};

export default PreviewContainer;
