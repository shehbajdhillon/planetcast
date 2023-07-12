import { Project, Transformation } from "@/types";
import {
  Box,
  HStack,
  Text,
  Button,
  useColorModeValue,
  Spinner,
} from "@chakra-ui/react";
import { useRouter } from "next/router";
import { useState } from "react";
import VideoPlayer from "../video_player";

interface ProjectCardProps {
  teamSlug: string;
  project: Project
};

const ProjectCard: React.FC<ProjectCardProps> = (props) => {

  const { project, teamSlug } = props;
  const transformations: Transformation[] = project?.transformations;
  const [transformationIdx, setTransformationIdx] = useState(0);

  const router = useRouter();

  return (
    <Box
      w={{ base: "330px", md: "360px" }}
      rounded={"lg"}
      _hover={{
        borderColor: useColorModeValue('gray.300', 'whiteAlpha.500'),
        bg: useColorModeValue('blackAlpha.50', 'whiteAlpha.200'),
      }}
      cursor={"pointer"}
    >
      <HStack pb="2px" pt="10px" rounded={"lg"}>
        <VideoPlayer src={transformations.length ? transformations?.[0].targetMedia : project.sourceMedia } style={{ borderRadius: "100px" }}/>
      </HStack>
      <Box p="5px" onClick={() => router.push(`/${teamSlug}/${project.id}`)}>
        <HStack p="5px">
          <Text
            textTransform="capitalize"
            fontWeight={700}
            fontSize={'lg'}
            letterSpacing={1.1}
            noOfLines={1}
          >
            {project.title}
          </Text>
        </HStack>
        <HStack px="5px" pb="5px">
          <Button
            borderWidth="1px"
            size={'xs'}
            textTransform="capitalize"
            fontWeight="medium"
            alignContent="right"
            pointerEvents={"none"}
          >
            {transformations.length ? transformations?.[0].targetLanguage : <Text>PROCESSING <Spinner size={"xs"} /></Text>}
          </Button>
        </HStack>
      </Box>
    </Box>
  );
};

export default ProjectCard;
