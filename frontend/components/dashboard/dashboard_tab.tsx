import { Box, Center, Divider, Grid, GridItem, Heading, useBreakpointValue, useColorModeValue } from "@chakra-ui/react";
import CastCard from "./cast_card";
import Image from "next/image";
import useWindowDimensions from "@/hooks/useWindowDimensions";

const DashboardTab: React.FC = () => {

  const dividerColor = useColorModeValue("gray.300", "whiteAlpha.300");
  const imageSize = useBreakpointValue({ base: 70 });

  const { height } = useWindowDimensions();

  return (
    <Box w="full" h="full" display={"flex"} flexDir={"column"}>
      <Box display={{ "lg": "none" }}>
        <Heading p="5px" fontSize={"3xl"}>Dashboard</Heading>
        <Divider borderColor={dividerColor} />
      </Box>
      <Center>
        <Grid
          px="100px"
          py={"30px"}
          placeItems={'center'}
          columnGap={'25px'}
          rowGap={'25px'}
          templateColumns={[
            'repeat(2, 1fr)',
            'repeat(2, 1fr)',
            'repeat(4, 1fr)',
            'repeat(6, 1fr)',
          ]}
          overflow={"auto"}
          maxH={{ base: (height as number) - 40, lg: height}}
        >
          <GridItem colSpan={2} w="full" h="full">
            <Box
              borderWidth={"1px"}
              w="full"
              h={"full"}
              p={6}
              rounded={"lg"}
              _hover={{
                borderColor: useColorModeValue('gray.300', 'whiteAlpha.500'),
                boxShadow: 'lg',
                bg: useColorModeValue('white', 'whiteAlpha.100'),
              }}
              cursor={"pointer"}
            >
              <Center h="full" w="full">
                <Image
                  src={useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg')}
                  width={imageSize}
                  height={100}
                  style={{ borderRadius: "20px" }}
                  alt='planet cast logo'
                />
              </Center>
            </Box>
          </GridItem>
          {new Array(20).fill(0).map(_ => (
            <GridItem colSpan={2}>
              <CastCard
                title="This Week In Startups"
                status="DRAFT"
                totalSteps={6}
                completedSteps={4}
              />
            </GridItem>
          ))}
        </Grid>
      </Center>
    </Box>
  );
};

export default DashboardTab;
