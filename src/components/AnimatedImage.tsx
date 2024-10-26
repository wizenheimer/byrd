"use client";

import type React from "react";
import { useEffect, useRef, useState } from "react";
import dynamic from "next/dynamic";

interface AnimatedImageProps {
	imageSrc: string;
	alt?: string;
	className?: string;
	imageClassName?: string;
	parallaxStrength?: number;
}

const AnimatedImage: React.FC<AnimatedImageProps> = ({
	imageSrc,
	alt,
	className = "",
	imageClassName = "",
	parallaxStrength = 0.1,
}) => {
	const containerRef = useRef<HTMLDivElement>(null);
	const [offset, setOffset] = useState(0);
	const [isVisible, setIsVisible] = useState(false);

	useEffect(() => {
		const observer = new IntersectionObserver(
			([entry]) => {
				if (entry.isIntersecting) {
					setIsVisible(true);
					if (containerRef.current) {
						observer.unobserve(containerRef.current);
					}
				}
			},
			{ threshold: 0.1 },
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

	// Prevent right-click
	const handleContextMenu = (e: React.MouseEvent) => {
		e.preventDefault();
	};

	return (
		<div
			ref={containerRef}
			className={`relative overflow-hidden ${className}`}
			onContextMenu={handleContextMenu}
			onDragStart={(e) => e.preventDefault()}
			draggable="false"
			style={
				{
					userSelect: "none",
					WebkitUserSelect: "none",
					MozUserSelect: "none",
					msUserSelect: "none",
				} as React.CSSProperties
			}
		>
			<img
				src={imageSrc}
				alt={alt || "Image"}
				className={`w-full object-cover transition-all duration-1000 ease-out ${imageClassName}`}
				style={{
					transform: `translateY(${offset}px) scale(${isVisible ? 1 : 1.05})`,
					opacity: isVisible ? 1 : 0,
					willChange: "transform, opacity",
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

// Export as a client-only component
export default dynamic(() => Promise.resolve(AnimatedImage), { ssr: false });
