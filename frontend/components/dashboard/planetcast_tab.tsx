import { Box, Heading, Divider, useColorModeValue } from "@chakra-ui/react";

const PlanetCastTab: React.FC = () => {

  const dividerColor = useColorModeValue("gray.300", "whiteAlpha.300");

  return (
    <Box w="full" h="full" display={"flex"} flexDir={"column"}>
      <Box display={{ "lg": "none" }}>
        <Heading p="5px" fontSize={"3xl"}>PlanetCast</Heading>
        <Divider borderColor={dividerColor} />
      </Box>
    </Box>
  );
};

export default PlanetCastTab;
