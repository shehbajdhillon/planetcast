import { Box, keyframes, BoxProps } from '@chakra-ui/react';
import React from 'react';

const scroll = keyframes`
  0% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(-50%);
  }
`;

interface MarqueeProps extends BoxProps {
  speed?: number;
}

const Marquee: React.FC<MarqueeProps> = ({ children, speed = 20, ...props }) => {
  // Convert children into an array and duplicate it
  const childrenArray = React.Children.toArray(children).concat(React.Children.toArray(children));

  return (
    <Box overflow="hidden" {...props}>
      <Box
        display="flex"
        whiteSpace="nowrap"
        animation={`${scroll} ${speed}s linear infinite`}
      >
        {childrenArray.map((child, index) => (
          <Box as="span" flexShrink="0" key={index}>
            {child}
          </Box>
        ))}
      </Box>
    </Box>
  );
};

export default Marquee;

