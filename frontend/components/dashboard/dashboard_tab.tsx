import { Box, Divider, Heading, useColorModeValue } from "@chakra-ui/react";

const DashboardTab: React.FC = () => {

  const dividerColor = useColorModeValue("gray.300", "whiteAlpha.300");

  return (
    <Box w="full" h="full" display={"flex"} flexDir={"column"}>
      <Box display={{ "lg": "none" }}>
        <Heading p="5px" fontSize={"3xl"}>Dashboard</Heading>
        <Divider borderColor={dividerColor} />
      </Box>
    </Box>
  );
};

export default DashboardTab;
