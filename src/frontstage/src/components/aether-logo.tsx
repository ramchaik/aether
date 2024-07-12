import React from "react";

interface AetherLogoProps {
  width?: number;
  height?: number;
}

const AetherLogo: React.FC<AetherLogoProps> = ({
  width = 200,
  height = 200,
}) => (
  <svg
    width={width}
    height={height}
    viewBox="0 0 200 200"
    xmlns="http://www.w3.org/2000/svg"
  >
    <defs>
      <linearGradient id="grad1" x1="0%" y1="0%" x2="100%" y2="100%">
        <stop offset="0%" style={{ stopColor: "#4A90E2", stopOpacity: 1 }} />
        <stop offset="100%" style={{ stopColor: "#5C2D91", stopOpacity: 1 }} />
      </linearGradient>
    </defs>

    <circle cx="100" cy="100" r="90" fill="url(#grad1)" />
    <path
      d="M60 160 L100 40 L140 160 Z"
      fill="none"
      stroke="white"
      strokeWidth="8"
    />
    <line x1="75" y1="120" x2="125" y2="120" stroke="white" strokeWidth="8" />
    <circle cx="50" cy="50" r="3" fill="white" />
    <circle cx="150" cy="50" r="4" fill="white" />
    <circle cx="100" cy="25" r="2" fill="white" />
    <circle cx="175" cy="100" r="3" fill="white" />
    <circle cx="25" cy="100" r="2" fill="white" />
  </svg>
);

export default AetherLogo;
