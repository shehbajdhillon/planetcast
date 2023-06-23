import {
  Box,
  Center,
  Grid,
  GridItem,
  useBreakpointValue,
  useColorModeValue,
  useDisclosure,
} from "@chakra-ui/react";
import CastCard from "./cast_card";
import Image from "next/image";
import NewCastModal from "../new_cast_modal";

interface DashboardTabProps {
  teamSlug: string;
};

const DashboardTab: React.FC<DashboardTabProps> = ({ teamSlug }) => {
  const imageSize = useBreakpointValue({ base: 70 });
  const { onOpen, isOpen, onClose } = useDisclosure();

  return (
    <Box w="full" h="full" display={"flex"} flexDir={"column"}>
      <NewCastModal onOpen={onOpen} isOpen={isOpen} onClose={onClose} />
      <Center>
        <Grid
          py={{ base: "100px" }}
          px={{ base: "35px", lg: "70px" }}
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
              maxW="400px"
              minW="270px"
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
                title="This Week In Startups Ep 1024"
                status="DONE"
                teamSlug={teamSlug}
                projectId={idx}
              />
            </GridItem>
          ))}
        </Grid>
      </Center>
    </Box>
  );
};

export default DashboardTab;
