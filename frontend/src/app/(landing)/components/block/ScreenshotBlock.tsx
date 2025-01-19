import AnimatedImage from "@/components/AnimatedImage";

type Props = {
	imageSrc: string;
};

const ScreenshotBlock = (props: Props) => {
	return (
		<div className="max-w-6xl mx-auto px-4 relative">
			<AnimatedImage imageSrc={props.imageSrc} />
		</div>
	);
};

export default ScreenshotBlock;
