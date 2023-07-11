import {
  Box,
  Center,
  Grid,
  GridItem,
  Text,
  useColorModeValue,
  useDisclosure,
} from "@chakra-ui/react";
import ProjectCard from "./project_card";
import Image from "next/image";
import NewProjectModal from "../new_project_modal";
import { Project } from "@/types";

interface DashboardTabProps {
  teamSlug: string;
  projects: Project[];
  refetch: () => void;
};

const LoadingBox: React.FC = () => {
  const bgColor = useColorModeValue("blackAlpha.50", "whiteAlpha.200")
  return (
    <Box
      w={{ base: '330px', md: '360px' }}
      h={{ base: '282px', md: '300px' }}
      rounded={"lg"}
      borderWidth={"1px"}
      borderColor={bgColor}
    >
      <Box py="10px" h={{ base: "185px", md: "202px" }} backgroundColor={bgColor} rounded={"lg"} />
      <Box m="10px" h="27px" w="100px" backgroundColor={bgColor} rounded={"lg"} />
      <Box m="10px" h="27px" w="50px" backgroundColor={bgColor} rounded={"lg"} />
    </Box>
  );
};

const LoadingGrid: React.FC = () => {
  return (
    <Grid
      py={{ base: "100px" }}
      px={{ base: "35px", lg: "70px" }}
      placeItems={'center'}
      rowGap={'60px'}
      columnGap={'60px'}
      templateColumns={{
        base: 'repeat(2, 1fr)',
        md: 'repeat(4, 1fr)',
        xl: 'repeat(6, 1fr)'
      }}
    >
      {new Array(8).fill(0).map((_, idx: number) => (
        <GridItem colSpan={2} key={idx}>
          <LoadingBox />
        </GridItem>
      ))}
    </Grid>
  );
};

const DashboardTab: React.FC<DashboardTabProps> = ({ teamSlug, projects, refetch }) => {
  const { onOpen, isOpen, onClose } = useDisclosure();
  const borderColor = useColorModeValue('gray.300', 'whiteAlpha.500');
  const bgColor = useColorModeValue('white', 'black');
  const imgSrc = useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg');

  return (
    <Box w="full" h="full" display={"flex"} flexDir={"column"}>
      <NewProjectModal onOpen={onOpen} isOpen={isOpen} onClose={onClose} refetch={refetch} teamSlug={teamSlug} />
      <Center>
      {!projects ? <LoadingGrid /> :
        <Grid
          py={{ base: "100px" }}
          px={{ base: "35px", lg: "70px" }}
          placeItems={'center'}
          rowGap={'60px'}
          columnGap={'60px'}
          templateColumns={{
            base: 'repeat(2, 1fr)',
            md: 'repeat(4, 1fr)',
            xl: 'repeat(6, 1fr)'
          }}
        >
          <GridItem colSpan={2}>
            <Box
              p={5}
              onClick={onOpen}
              rounded={"lg"}
              _hover={{
                borderColor: borderColor,
                borderWidth: "1px",
                boxShadow: 'lg',
                bg: bgColor,
              }}
              cursor={"pointer"}
            >
              <Center h="full" flexDirection={"column"}>
                <Image
                  src={imgSrc}
                  width={70}
                  height={100}
                  style={{ borderRadius: "20px" }}
                  alt='planet cast logo'
                />
                <Text>New Project</Text>
              </Center>
            </Box>
          </GridItem>

          {projects?.map((project: Project, idx: number) => (
            <GridItem colSpan={2} key={idx}>
              <ProjectCard
                teamSlug={teamSlug}
                project={project}
              />
            </GridItem>
          ))}
        </Grid>
      }
      </Center>
    </Box>
  );
};

export default DashboardTab;
