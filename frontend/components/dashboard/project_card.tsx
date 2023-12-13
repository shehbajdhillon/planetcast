import { Project } from "@/types";
import {
  Box,
  HStack,
  Text,
  Button,
  useColorModeValue,
  Spinner,
  Spacer,
  IconButton,
  Center,
  Progress,
  VStack,
} from "@chakra-ui/react";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import VideoPlayer from "../video_player";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { gql, useLazyQuery } from "@apollo/client";

interface ProjectCardProps {
  teamSlug: string;
  project: Project
};

const GET_PROJECT_DATA = gql`
  query GetProjectData($teamSlug: String!, $projectId: Int64!) {
    getTeamById(teamSlug: $teamSlug) {
      projects(projectId: $projectId) {
        id
        title
        sourceMedia
        transformations {
          id
          targetMedia
          targetLanguage
          status
          progress
        }
      }
    }
  }
`;


interface LoadingBoxProps {
  progress: number;
  status: string;
}

const VideoProcessingBox: React.FC = () => {
  return (
    <Box
      w={{ base: "330px", md: "360px" }}
      h={{ base: "185.63px", md:"202.5px" }}
    >
      <VStack mt={"55px"} mx="5px" w="full">
        <Spinner size={'xl'} />
        <Center pt="10px" fontSize={'xl'}>{"Processing Upload"}</Center>
      </VStack>
    </Box>
  );
}

const LoadingBox: React.FC<LoadingBoxProps> = ({ status, progress }) => {
  const bgColor = useColorModeValue("gray.200", "gray.800");
  return (
    <Box
      w={{ base: "330px", md: "360px" }}
      h={{ base: "185.63px", md:"202.5px" }}
    >
        { status !== "error" ?
          <Box mt={"110px"} px="5px" w="full">
            <Progress value={progress} hasStripe size="md" isAnimated={true} rounded={"sm"} backgroundColor={bgColor} />
            <Center pt="10px">{progress + "%"}</Center>
          </Box>
          :
          <Box mt={"110px"} px="5px" w="full">
            <Progress value={100} size="md" rounded={"sm"} backgroundColor={bgColor} colorScheme="red" />
            <Center pt="10px">
              <Text
                noOfLines={1}
              >
                {"An error occured while processing your dubbing"}
              </Text>
            </Center>
          </Box>
        }
    </Box>
  );
}


const ProjectCard: React.FC<ProjectCardProps> = (props) => {

  const { project: currentProject, teamSlug } = props;

  const [project, setProject] = useState<Project>(currentProject);
  const [transformationIdx, setTransformationIdx] = useState(0);

  const transformations = project?.transformations;
  const transformation = transformations && transformations[transformationIdx];
  const currentStatus = transformation && transformation?.status;


  const isProcessing = transformations?.length === 0 || currentStatus !== "complete"

  const [getProjectData, { data }]
    = useLazyQuery(GET_PROJECT_DATA, {
        variables: { teamSlug, projectId: project.id },
        fetchPolicy: 'no-cache',
        pollInterval: isProcessing ? 10000 : 0
      })

  useEffect(() => {
    const newProjectData = data?.getTeamById.projects[0];
    if (newProjectData) {
      setProject(newProjectData);
    }
  }, [data]);

  useEffect(() => {
    if (isProcessing) {
      getProjectData();
    }
  }, [isProcessing]);

  const router = useRouter();

  const goNext = () => setTransformationIdx(curr => curr + 1)
  const goBack = () => setTransformationIdx(curr => curr - 1)

  const borderColor = useColorModeValue('gray.300', 'whiteAlpha.500');
  const bg = useColorModeValue('blackAlpha.50', 'whiteAlpha.200');

  return (
    project &&
    <Box
      w={{ base: "330px", md: "360px" }}
      rounded={"lg"}
      _hover={{
        borderColor: borderColor,
        bg: bg
      }}
      cursor={"pointer"}
    >
      <HStack pb="2px" pt="10px" rounded={"lg"}>
        {
          (isProcessing && transformations?.length) ? <LoadingBox progress={transformations?.[transformationIdx].progress} status={currentStatus} /> :
          project.sourceMedia === "" ? <VideoProcessingBox /> : <VideoPlayer src={transformations.length ? transformations?.[transformationIdx].targetMedia : project.sourceMedia } style={{ borderRadius: "100px" }}/>
        }
      </HStack>
      <Box p="5px" onClick={() => router.push(`/dashboard/${teamSlug}/project/${project.id}`)}>
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
            <HStack spacing={"4px"}>
              <Text>{!transformations.length ? "PROCESSING" : transformations?.[transformationIdx]?.targetLanguage}</Text>
              { isProcessing && <Spinner size={"xs"} /> }
            </HStack>
          </Button>
          <Spacer />
          <IconButton
            icon={<ChevronLeft />}
            aria-label="choose previous dubbing"
            variant={"outline"}
            isDisabled={transformationIdx - 1 < 0}
            onClick={(e) => {
              e.stopPropagation();
              goBack();
            }}
          />
          <IconButton
            icon={<ChevronRight />}
            aria-label="choose next dubbing"
            variant={"outline"}
            isDisabled={transformationIdx + 1 >= transformations.length}
            onClick={(e) => {
              e.stopPropagation();
              goNext();
            }}
          />
        </HStack>
      </Box>
    </Box>
  );
};

export default ProjectCard;
