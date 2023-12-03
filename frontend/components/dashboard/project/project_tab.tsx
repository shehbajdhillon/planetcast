import { useVideoSeekStore } from "@/stores/video_seek_store";
import { Project, Segment } from "@/types";
import { gql, useLazyQuery, useMutation } from "@apollo/client";
import {
  AspectRatio,
  Box,
  Button,
  Center,
  Checkbox,
  Grid,
  GridItem,
  HStack,
  IconButton,
  Progress,
  SkeletonText,
  Spinner,
  Text,
  VStack,
  useClipboard,
  useColorModeValue,
  useDisclosure
} from "@chakra-ui/react";
import { Clipboard } from "lucide-react";
import { useEffect, useState } from "react";
import Link from "next/link";
import VideoPlayer from "@/components/video_player";
import NewTransformationModel from "@/components/dashboard/project/new_transformation_modal";
import { DownloadIcon, TrashIcon } from "lucide-react";
import { formatTime } from "@/utils";
import SingleActionModal from "@/components/single_action_modal";

const DELETE_TRANSFORMATION = gql`
  mutation DeleteTransformation($transformationId: Int64!) {
    deleteTransformation(transformationId: $transformationId) {
      id
    }
  }
`;

const GET_CURRENT_PROJECT = gql`
  query GetCurrentProject($teamSlug: String!, $projectId: Int64!) {
    getTeamById(teamSlug: $teamSlug) {
      projects(projectId: $projectId) {
        id
        title
        sourceMedia
        dubbingCreditsRequired
        transformations {
          id
          targetMedia
          targetLanguage
          transcript
          isSource
          status
          progress
        }
      }
    }
  }
`;

interface ProjectTabProps {
  project: Project;
  teamSlug: string;
};



const VideoProcessingBox: React.FC = () => {
  return (
    <AspectRatio ratio={16/9}>
      <Box
        w={"full"}
        maxW={"1280px"}
      >
        <VStack mt={"0px"} mx="5px" w="full">
          <Spinner size={'xl'} />
          <Center pt="10px" fontSize={'xl'}>{"Processing Upload"}</Center>
        </VStack>
      </Box>
    </AspectRatio>
  )
};


interface LoadingBoxProps {
  progress: number;
  status: string;
}

const LoadingBox: React.FC<LoadingBoxProps> = ({ status, progress }) => {
  const bgColor = useColorModeValue("gray.200", "gray.800");
  return (
    <AspectRatio ratio={16/9}>
      <Box
        w={"full"}
        maxW={"1280px"}
      >
          { status !== "error" ?
            <Box mt={"0px"} mx="5px" w="full">
              <Progress value={progress} hasStripe size="md" isAnimated={true} rounded={"sm"} backgroundColor={bgColor} />
              <Center pt="10px">{progress + "%"}</Center>
            </Box>
            :
            <Box mt={"0px"} mx="5px" w="full">
              <Progress value={100} size="md" rounded={"sm"} backgroundColor={bgColor} colorScheme="red" />
              <Center pt="10px">{"An error occured while processing your dubbing, please delete and try again. All credits used for this dubbing have been refunded."}</Center>
            </Box>
          }
      </Box>
    </AspectRatio>
  );
}


interface TranscriptViewProps {
  segments: Segment[];
};

const TranscriptView: React.FC<TranscriptViewProps> = ({ segments }) => {

  const currentSeek = useVideoSeekStore((state) => state.currentSeek);
  const [autoScroll, setAutoScroll] = useState(true);

  const parentId = "transcript-parent-view"
  const transcriptMessageId = (segment: Segment) => `transcript-message-${segment.id}`

  const highlight = (segment: Segment) => {
    if (segment.start <= currentSeek && currentSeek <= segment.end) {
      if (autoScroll) {
        const topPosition = document.getElementById(transcriptMessageId(segment))?.offsetTop;
        const parentDiv = document.getElementById(parentId);
        if (parentDiv && topPosition) parentDiv.scrollTop = topPosition;
      }
      return true;
    };
    return false;
  };


  useEffect(() => {
    const parentDiv = document.getElementById(parentId);
    if (parentDiv) parentDiv.scrollTop = 0;
  }, []);

  return (
    <VStack h="full" p="10px">
      <Checkbox
        w="full"
        pl="1px"
        colorScheme="gray"
        isChecked={autoScroll}
        onChange={(e) => setAutoScroll(e.target.checked)}
      >
        Auto Scroll
      </Checkbox>
      <VStack overflow={"scroll"} h="full" id={parentId} position={"relative"} w="full">
        {segments?.map((segment: Segment, idx: number) => (
          <>
            <MessageView
              segment={segment}
              key={idx}
              htmlId={transcriptMessageId(segment)}
              highlight={highlight(segment)}
            />
          </>
        ))}
      </VStack>
    </VStack>
  );
};


interface MessageViewProps {
  segment: Segment;
  htmlId: string;
  highlight: boolean;
};

const MessageView: React.FC<MessageViewProps> = ({ segment, htmlId, highlight }) => {

  const bgColorHighlight = useColorModeValue("black", "white");
  const textColor = useColorModeValue("white", "black");

  const { onCopy } = useClipboard(segment.text);

  return (
    <Box
      rounded="10px"
      whiteSpace={'normal'}
      height="auto"
      blockSize={'auto'}
      w="full"
      justifyContent="left"
      borderWidth={"1px"}
      textColor={highlight ? textColor : 'inherit'}
      bgColor={highlight ? bgColorHighlight : 'inherit'}
      id={htmlId}
      p={2}
    >
      <HStack>
        <IconButton
          aria-label="copy transcript message"
          icon={<Clipboard size={"20px"} />}
          size={"0px"}
          onClick={onCopy}
          variant={"unstyled"}
        />
        <Text>{formatTime(segment.start)} - {formatTime(segment.end)}</Text>
      </HStack>
      <Text
        textAlign={"left"}
        fontWeight={"semibold"}
      >
        { segment.text.trim() }
      </Text>
    </Box>
  );
};

interface LoadingTranscriptView {
  transcribing: boolean;
};

export const LoadingTranscriptView: React.FC<LoadingTranscriptView> = ({ transcribing }) => {
  return (
    <VStack h="full" p="10px">
      <VStack overflow={"hidden"} h="full" position={"relative"} w="full">
        { transcribing &&
          <HStack>
            <Text>Video is currently being transcribed. This may take 5-15 minutes.</Text>
            <Spinner />
          </HStack>
        }
        {new Array(7).fill(0).map((_, idx: number) => (
          <LoadingMessageView key={idx} />
        ))}
      </VStack>
    </VStack>
  );
}

const LoadingMessageView: React.FC = () => {

  return (
    <Box
      rounded="10px"
      whiteSpace={'normal'}
      height="auto"
      blockSize={'auto'}
      w="full"
      justifyContent="left"
      borderWidth={"1px"}
      p={2}
    >
      <SkeletonText noOfLines={4} h={"66px"} rounded={"10px"} />
    </Box>
  );

};


const ProjectTab: React.FC<ProjectTabProps> = (props) => {

  const { project: currentProject, teamSlug } = props;

  const [project, setProject] = useState<Project>(currentProject);
  const [transformationIdx, setTransformationIdx] = useState(0);

  const transformations = project && project?.transformations;
  const transformation = transformations && transformations[transformationIdx];
  const currentStatus = transformation && transformation?.status;
  const parseTranscript = transformation && transformation.transcript && JSON.parse(transformation.transcript)

  const isProcessing = transformations?.length === 0 || currentStatus !== "complete"

  const [getProjectData, { data, refetch }]
    = useLazyQuery(GET_CURRENT_PROJECT, { variables: { teamSlug, projectId: project?.id }, pollInterval: !isProcessing ? 0 : 10000, fetchPolicy: 'no-cache' });

  const [deleteTransformation, { loading: deleteTfnLoading } ] = useMutation(DELETE_TRANSFORMATION);

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

  const { isOpen, onClose, onOpen } = useDisclosure();

  const setCurrentSeek = useVideoSeekStore((state) => state.setCurrentSeek);
  const onTimeUpdate = (time: number) => setCurrentSeek(time);

  const deleteDubbing = async (tfnId: number) => {
    const res = await deleteTransformation({ variables: { transformationId: tfnId } });
    if (res) {
      setTransformationIdx(0);
      refetch();
    }
  }

  useEffect(() => {
    setCurrentSeek(0);
  }, [setCurrentSeek]);

  const buttonColor = useColorModeValue("white", "black");
  const buttonBg = useColorModeValue("black", "white");

  return (
    project &&
    <Box
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      w={"full"}
    >
      <Box w="full" maxW={"1920px"}>

        <Grid
          templateAreas={{
            base: `
              "video"
              "transcript"
            `,
            lg: `"video transcript"`
          }}
          gridTemplateColumns={{ base: "1fr", lg: "2fr 1fr" }}
          w="full"
          h={"full"}
          gap={"10px"}
        >

          <GridItem area={'video'} h="full" w="full" rounded={"lg"} maxW={"1280px"}>
            {
              (isProcessing && transformation) ? <LoadingBox status={transformation.status} progress={transformation.progress} /> :
              project.sourceMedia === "" ? <VideoProcessingBox /> : <VideoPlayer src={transformation ? transformation?.targetMedia : project.sourceMedia } onTimeUpdate={onTimeUpdate} />
            }
            <HStack display="flex" flexWrap={"wrap"} overflow={"auto"} spacing={"10px"} pt="10px" hidden={!transformations?.length}>
              { transformations?.map((t: any, idx: number) => (
                <Button
                  key={idx}
                  onClick={() => setTransformationIdx(idx)}
                  variant={idx === transformationIdx ? "solid" : "outline"}
                  pointerEvents={idx === transformationIdx ? "none" : "auto"}
                  background={idx === transformationIdx ? buttonBg : '' }
                  color={idx === transformationIdx ? buttonColor : '' }
                >
                  {t?.targetLanguage}
                </Button>
              ))}
              <NewTransformationModel project={project} refetch={refetch} />
              { transformation && <Link href={transformation?.targetMedia}><Button leftIcon={<DownloadIcon />} variant={"outline"}>Download</Button></Link> }
              <SingleActionModal
                heading={"Delete Dubbing"}
                action={() => deleteDubbing(transformations[transformationIdx].id)}
                loading={deleteTfnLoading}
                isOpen={isOpen}
                onClose={onClose}
              >
                Are you sure you want to delete this dubbing? This action is irreversible.
              </SingleActionModal>
              {!transformations[transformationIdx]?.isSource && <Button onClick={onOpen} leftIcon={<TrashIcon />} variant={"outline"}>Delete Dubbing</Button>}
            </HStack>
          </GridItem>

          <GridItem area={'transcript'} h="full" w="full" borderWidth={"1px"} rounded="lg" maxH={"596px"}>

          { !parseTranscript ?
            <LoadingTranscriptView transcribing={transformations?.length <= 0 || transformations === undefined} />
            :
            <TranscriptView segments={parseTranscript?.segments} />
          }
          </GridItem>

        </Grid>
      </Box>
    </Box>
  );
};

export default ProjectTab;
