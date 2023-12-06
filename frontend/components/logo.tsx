import React from 'react';

interface LogoProps {
  height: number;
  width: number;
};

export const LightModeGradientLogo: React.FC<LogoProps> = ({ height, width }) => (
  <svg width={width} height={height} viewBox="0 0 600 600" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M0 50C0 22.3858 22.3858 0 50 0H550C577.614 0 600 22.3858 600 50V550C600 577.614 577.614 600 550 600H50C22.3858 600 0 577.614 0 550V50Z" fill="white"/>
    <ellipse cx="312.936" cy="47.9391" rx="312.936" ry="47.9391" transform="matrix(0.72336 -0.690472 0.659889 0.751363 40.8 480.147)" fill="url(#paint0_linear_4_16)"/>
    <ellipse cx="270.209" cy="41.3937" rx="270.209" ry="41.3937" transform="matrix(0.72336 -0.690472 0.659889 0.751363 74.3148 456.084)" fill="white"/>
    <circle cx="298.8" cy="298.8" r="174" fill="url(#paint1_linear_4_16)"/>
    <path d="M510.143 131.807C495.963 154.652 473.586 183.849 444.958 216.859C416.331 249.869 382.319 285.692 345.883 321.212C309.446 356.732 271.686 390.873 235.889 420.667C200.092 450.461 167.341 475.005 140.484 492.164L171.224 454.869C192.866 441.042 219.258 421.264 248.105 397.255C276.951 373.246 307.379 345.734 336.741 317.111C366.102 288.488 393.51 259.621 416.579 233.021C439.648 206.421 457.68 182.893 469.106 164.483L510.143 131.807Z" fill="white"/>
    <defs>
      <linearGradient id="paint0_linear_4_16" x1="312.936" y1="0" x2="312.936" y2="95.8783" gradientUnits="userSpaceOnUse">
        <stop stop-color="#007CF0"/>
        <stop offset="0.71875" stop-color="#01DFD8" stop-opacity="0"/>
      </linearGradient>
      <linearGradient id="paint1_linear_4_16" x1="298.8" y1="124.8" x2="298.8" y2="472.8" gradientUnits="userSpaceOnUse">
        <stop stop-color="#007CF0"/>
        <stop offset="1" stop-color="#01DFD8"/>
      </linearGradient>
    </defs>
  </svg>
);


export const DarkModeGradientLogo: React.FC<LogoProps> = ({ height, width }) => (
  <svg width={width} height={height} viewBox="0 0 600 600" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M0 50C0 22.3858 22.3858 0 50 0H550C577.614 0 600 22.3858 600 50V550C600 577.614 577.614 600 550 600H50C22.3858 600 0 577.614 0 550V50Z" fill="black"/>
    <ellipse cx="312.936" cy="47.9391" rx="312.936" ry="47.9391" transform="matrix(0.72336 -0.690472 0.659889 0.751363 40.8 480.147)" fill="url(#paint0_linear_7_8)"/>
    <ellipse cx="270.209" cy="41.3937" rx="270.209" ry="41.3937" transform="matrix(0.72336 -0.690472 0.659889 0.751363 74.3148 456.084)" fill="black"/>
    <circle cx="298.8" cy="298.8" r="174" fill="url(#paint1_linear_7_8)"/>
    <path d="M510.143 131.807C495.963 154.652 473.586 183.849 444.958 216.859C416.331 249.869 382.319 285.691 345.883 321.211C309.446 356.732 271.686 390.873 235.889 420.667C200.092 450.461 167.341 475.005 140.484 492.164L171.224 454.869C192.866 441.042 219.258 421.264 248.105 397.255C276.951 373.246 307.379 345.734 336.741 317.111C366.102 288.488 393.51 259.621 416.579 233.021C439.648 206.42 457.68 182.893 469.106 164.483L510.143 131.807Z" fill="black"/>
    <defs>
      <linearGradient id="paint0_linear_7_8" x1="312.936" y1="0" x2="312.936" y2="95.8783" gradientUnits="userSpaceOnUse">
        <stop stop-color="#007CF0"/>
        <stop offset="0.71875" stop-color="#01DFD8" stop-opacity="0"/>
      </linearGradient>
      <linearGradient id="paint1_linear_7_8" x1="298.8" y1="124.8" x2="298.8" y2="472.8" gradientUnits="userSpaceOnUse">
        <stop stop-color="#007CF0"/>
        <stop offset="1" stop-color="#01DFD8"/>
      </linearGradient>
    </defs>
  </svg>
);
