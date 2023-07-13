import { Project, Transformation } from "@/types";
import {
  Box,
  HStack,
  Text,
  Button,
  useColorModeValue,
  Spinner,
  Spacer,
  IconButton,
} from "@chakra-ui/react";
import { useRouter } from "next/router";
import { useState } from "react";
import VideoPlayer from "../video_player";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface ProjectCardProps {
  teamSlug: string;
  project: Project
};

const ProjectCard: React.FC<ProjectCardProps> = (props) => {

  const { project, teamSlug } = props;
  const transformations: Transformation[] = project?.transformations;
  const [transformationIdx, setTransformationIdx] = useState(0);

  const router = useRouter();

  const goNext = () => setTransformationIdx(curr => curr + 1)
  const goBack = () => setTransformationIdx(curr => curr - 1)

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
        <VideoPlayer src={transformations.length ? transformations?.[transformationIdx].targetMedia : project.sourceMedia } style={{ borderRadius: "100px" }}/>
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
            {transformations.length ? transformations?.[transformationIdx].targetLanguage : <Text>PROCESSING <Spinner size={"xs"} /></Text>}
          </Button>
          <Spacer />
          <IconButton
            icon={<ChevronLeft />}
            aria-label="choose previous dubbing"
            variant={"outline"}
            isDisabled={transformationIdx - 1 < 0}
            onClick={(e) => {
              e.stopPropagation();
              setTransformationIdx(curr => curr - 1)
            }}
          />
          <IconButton
            icon={<ChevronRight />}
            aria-label="choose next dubbing"
            variant={"outline"}
            isDisabled={transformationIdx + 1 >= transformations.length}
            onClick={(e) => {
              e.stopPropagation();
              setTransformationIdx(curr => curr + 1)
            }}
          />
        </HStack>
      </Box>
    </Box>
  );
};

export default ProjectCard;
