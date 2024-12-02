// src/components/AnimatedImage.tsx
"use client";

import dynamic from "next/dynamic";
import type React from "react";
import { useEffect, useRef, useState } from "react";

interface RevealImageProps {
  imageSrc: string;
  alt?: string;
  className?: string;
  imageClassName?: string;
}

const RevealImage: React.FC<RevealImageProps> = ({
  imageSrc,
  alt,
  className = "",
  imageClassName = "",
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
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
        className={`w-full object-cover transition-all duration-700 ease-out ${imageClassName}`}
        style={{
          transform: `scale(${isVisible ? 1 : 1.05})`,
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

export default dynamic(() => Promise.resolve(RevealImage), { ssr: false });
