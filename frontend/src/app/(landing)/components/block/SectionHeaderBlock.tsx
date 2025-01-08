// src/components/block/SectionHeaderBlock.tsx
import React from "react";

type SectionHeaderProps = {
  title: string;
  subtitle?: string;
  className?: string;
  containerClassName?: string;
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

const SectionHeader = ({
  title,
  subtitle,
  className = "",
  containerClassName = "",
}: SectionHeaderProps) => {
  return (
    <div
      className={`max-w-7xl mx-auto px-4 pt-12 md:pt-16 lg:pt-20 pb-20 md:pb-32 lg:pb-40 ${containerClassName}`}
    >
      <div className="text-center max-w-4xl mx-auto">
        <h1 className={`text-5xl font-bold tracking-tight mb-6 ${className}`}>
          <MultiLineText text={title} />
        </h1>
        {subtitle && (
          <p className="text-lg text-gray-600 mt-8">
            <MultiLineText text={subtitle} />
          </p>
        )}
      </div>
    </div>
  );
};

export default SectionHeader;
