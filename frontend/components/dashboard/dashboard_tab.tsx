import { Box, Center, Divider, Grid, GridItem, Heading, useBreakpointValue, useColorModeValue, useDisclosure } from "@chakra-ui/react";
import CastCard from "./cast_card";
import Image from "next/image";
import NewCastModal from "../new_cast_modal";
import Navbar from "./navbar";

const DashboardTab: React.FC = () => {

  const dividerColor = useColorModeValue("gray.300", "whiteAlpha.300");
  const imageSize = useBreakpointValue({ base: 70 });

  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <Box w="full" h="full" display={"flex"} flexDir={"column"}>
      <NewCastModal onOpen={onOpen} isOpen={isOpen} onClose={onClose} />
      <Box display={{ "lg": "none" }}>
        <Heading p="5px" fontSize={"3xl"}>Dashboard</Heading>
        <Divider borderColor={dividerColor} />
      </Box>
      <Box position={"fixed"} top={0} left={0} w="full" p="10px" backgroundColor={useColorModeValue("white", "black")} zIndex={1000}>
        <Navbar />
      </Box>
      <Center>
        <Grid
          py={{ base: "50px", md: "50px", lg: "100px" }}
          placeItems={'center'}
          columnGap={'25px'}
          rowGap={'25px'}
          templateColumns={[
            'repeat(2, 1fr)',
            'repeat(2, 1fr)',
            'repeat(4, 1fr)',
            'repeat(6, 1fr)',
          ]}
        >
          <GridItem colSpan={2} w="full" h="full">
            <Box
              onClick={onOpen}
              borderWidth={"1px"}
              w="full"
              h={"full"}
              p={6}
              rounded={"lg"}
              _hover={{
                borderColor: useColorModeValue('gray.300', 'whiteAlpha.500'),
                boxShadow: 'lg',
                bg: useColorModeValue('white', 'black'),
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
          {new Array(20).fill(0).map((_, idx) => (
            <GridItem colSpan={2} key={idx}>
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
