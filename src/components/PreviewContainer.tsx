"use client";

import type React from "react";
import { useEffect, useRef, useState } from "react";
import Image from "next/image";

interface PreviewContainerProps {
	imageSrc: string;
	caption?: string;
	variant?: "default" | "compact" | "full" | "fitted";
	position?: "center" | "right" | "left";
	imageSize?: "contain" | "cover";
	className?: string;
	imageClassName?: string;
	parallaxStrength?: number;
}

const PreviewContainer: React.FC<PreviewContainerProps> = ({
	imageSrc,
	caption,
	variant = "default",
	position = "right",
	imageSize = "contain",
	className = "",
	imageClassName = "",
	parallaxStrength = 0.05,
}) => {
	const containerRef = useRef<HTMLDivElement>(null);
	const [offset, setOffset] = useState(0);
	const [isVisible, setIsVisible] = useState(false);

	useEffect(() => {
		const observer = new IntersectionObserver(
			([entry]) => {
				if (entry.isIntersecting) {
					setIsVisible(true);
					// Once revealed, we can stop observing
					if (containerRef.current) {
						observer.unobserve(containerRef.current);
					}
				}
			},
			{
				threshold: 0.3,
			},
		);

		if (containerRef.current) {
			observer.observe(containerRef.current);
		}

		return () => observer.disconnect();
	}, []);

	useEffect(() => {
		const handleScroll = () => {
			if (!containerRef.current) return;

			const rect = containerRef.current.getBoundingClientRect();
			const containerCenter = rect.top + rect.height / 2;
			const viewportCenter = window.innerHeight / 2;
			const distanceFromCenter = containerCenter - viewportCenter;

			setOffset(distanceFromCenter * parallaxStrength);
		};

		window.addEventListener("scroll", handleScroll);
		handleScroll();

		return () => window.removeEventListener("scroll", handleScroll);
	}, [parallaxStrength]);

	const containerStyles = {
		default: "min-h-[800px]",
		compact: "min-h-[500px]",
		full: "min-h-screen",
		fitted: "h-auto",
	};

	const positionStyles = {
		center: "inset-0 m-auto",
		right: "-right-[15%] -bottom-[20%]",
		left: "-left-[15%] -bottom-[20%]",
	};

	const imageSizeStyles = {
		contain: "object-contain",
		cover: "object-cover",
	};

	const getImageClassName = () => {
		if (variant === "fitted") {
			return "relative w-full h-full object-cover";
		}
		return `
      absolute
      ${positionStyles[position]}
      ${imageSizeStyles[imageSize]}
      rounded-lg
      ${variant === "default" ? "h-[900px]" : "h-full w-full"}
      ${imageClassName}
      transition-all duration-1000 ease-out
    `;
	};

	return (
		<div className={`relative ${className}`} ref={containerRef}>
			<div
				className={`relative bg-gray-100 rounded-md p-2 overflow-hidden ${containerStyles[variant]} transition-opacity duration-1000 ease-out`}
				style={{
					opacity: isVisible ? 1 : 0,
				}}
			>
				<div
					className={`rounded-lg overflow-hidden bg-white ${
						variant === "fitted" ? "h-auto" : "w-full h-full"
					}`}
				>
					<Image
						src={imageSrc}
						alt={caption || "Preview image"}
						className={getImageClassName()}
						style={{
							transform: `translateY(${offset}px) scale(${isVisible ? 1 : 1.05})`,
							opacity: isVisible ? 1 : 0,
							willChange: "transform, opacity",
						}}
					/>
				</div>
			</div>
			{caption && (
				<div
					className="relative mt-4 text-sm text-gray-600 px-1 transition-all duration-1000 ease-out"
					style={{
						opacity: isVisible ? 1 : 0,
						transform: `translateY(${isVisible ? 0 : "20px"})`,
					}}
				>
					{caption}
				</div>
			)}
		</div>
	);
};

export default PreviewContainer;
