import {
  Box,
  Center,
  Grid,
  GridItem,
  Text,
  useColorModeValue,
  useDisclosure,
} from "@chakra-ui/react";
import ProjectCard from "./cast_card";
import Image from "next/image";
import NewCastModal from "../new_cast_modal";
import { Project } from "@/types";

interface DashboardTabProps {
  teamSlug: string;
  projects: Project[];
  refetch: () => void;
};

const DashboardTab: React.FC<DashboardTabProps> = ({ teamSlug, projects, refetch }) => {
  const { onOpen, isOpen, onClose } = useDisclosure();
  return (
    <Box w="full" h="full" display={"flex"} flexDir={"column"}>
      <NewCastModal onOpen={onOpen} isOpen={isOpen} onClose={onClose} refetch={refetch} teamSlug={teamSlug} />
      <Center>
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

        <GridItem colSpan={2} h="full">
          <Box
            onClick={onOpen}
            borderWidth={"1px"}
            w={{ base: "330px", md: "360px" }}
            h={"full"}
            rounded={"lg"}
            _hover={{
              borderColor: useColorModeValue('gray.300', 'whiteAlpha.500'),
              boxShadow: 'lg',
              bg: useColorModeValue('white', 'black'),
            }}
            cursor={"pointer"}
          >
            <Center h="full" flexDirection={"column"}>
              <Image
                src={useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg')}
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
      </Center>
    </Box>
  );
};

export default DashboardTab;
