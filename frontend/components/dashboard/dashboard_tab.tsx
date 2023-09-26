import {
  Box,
  Center,
  Grid,
  GridItem,
  Skeleton,
  useColorModeValue,
} from "@chakra-ui/react";
import ProjectCard from "./project_card";
import NewProjectModal from "../new_project_modal";
import { Project } from "@/types";

interface DashboardTabProps {
  teamSlug: string;
  projects: Project[];
  refetch: () => void;
};

const LoadingBox: React.FC = () => {
  const bgColor = useColorModeValue("blackAlpha.300", "whiteAlpha.300")
  return (
    <Box
      w={{ base: '330px', md: '360px' }}
      h={{ base: '282px', md: '300px' }}
      rounded={"lg"}
      borderColor={bgColor}
    >
      <Skeleton py="10px" h={{ base: "185px", md: "202px" }} rounded={"lg"} />
      <Skeleton m="10px" h="27px" w="100px" rounded={"lg"} />
      <Skeleton m="10px" h="27px" w="50px" rounded={"lg"} />
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
  return (
    <Box w="full" h="full" display={"flex"} flexDir={"column"}>
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
          <GridItem colSpan={2} h="full" pt="10px">
            <NewProjectModal refetch={refetch} teamSlug={teamSlug} />
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
