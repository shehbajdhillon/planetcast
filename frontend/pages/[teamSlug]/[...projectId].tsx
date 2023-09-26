import Navbar from "@/components/dashboard/navbar";
import { Project, Segment, Team } from "@/types";
import { gql, useLazyQuery, useMutation, useQuery } from "@apollo/client";
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
  useDisclosure,
} from "@chakra-ui/react";
import { GetServerSideProps, NextPage } from "next";

import Head from "next/head";

import {
  Tabs,
  TabPanel,
  TabPanels,
  Tab,
  TabList
} from '@chakra-ui/react';
import { useRouter } from "next/router";
import VideoPlayer from "@/components/video_player";
import { useEffect, useState } from "react";
import { useVideoSeekStore } from "@/stores/video_seek_store";
import { formatTime } from "@/utils";
import NewTransformationModel from "@/components/new_transformation_modal";
import { Clipboard, TrashIcon } from "lucide-react";
import SingleActionModal from "@/components/single_action_modal";


const DELETE_PROJECT = gql`
  mutation DeleteProject($projectId: Int64!) {
    deleteProject(projectId: $projectId) {
      id
    }
  }
`;

interface SettingsTabProps {
  projectId: number;
  teamSlug: string;
};

const SettingsTab: React.FC<SettingsTabProps> = ({ projectId, teamSlug }) => {

  const [deleteProjectMutation, { loading }] = useMutation(DELETE_PROJECT);

  const router = useRouter();

  const deleteProject = async () => {
    const res = await deleteProjectMutation({ variables: { projectId } });
    if (res) router.push(`/${teamSlug}`);
  };

  const { isOpen, onClose, onOpen } = useDisclosure();

  return (
    <Box
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      w={"full"}
    >
      <Box w="full" maxW={"1920px"}>
        <Center>
          <SingleActionModal
            heading={"Delete Project"}
            body={`Are you sure you want to delete this Project? This will delete the original video and all the dubbings generated. This action is irreversible.`}
            action={() => deleteProject()}
            loading={loading}
            isOpen={isOpen}
            onClose={onClose}
          />
          <Button colorScheme="red" onClick={onOpen}>
            Delete Project
          </Button>
        </Center>
      </Box>
    </Box>
  );
};

interface ProjectTabProps {
  project: Project;
  teamSlug: string;
};

const GET_CURRENT_PROJECT = gql`
  query GetCurrentProject($teamSlug: String!, $projectId: Int64!) {
    getTeamById(teamSlug: $teamSlug) {
      projects(projectId: $projectId) {
        id
        title
        sourceMedia
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

const DELETE_TRANSFORMATION = gql`
  mutation DeleteTransformation($transformationId: Int64!) {
    deleteTransformation(transformationId: $transformationId) {
      id
    }
  }
`;

interface LoadingBoxProps {
  progress: number;
}

const LoadingBox: React.FC<LoadingBoxProps> = ({ progress }) => {
  return (
    <AspectRatio ratio={16/9}>
      <Box
        w={"full"}
        maxW={"1280px"}
      >
          <Box mt={"0px"} mx="5px" w="full">
            <Progress value={progress} hasStripe size="md" isAnimated={true} rounded={"sm"} backgroundColor={"gray.800"} />
            <Center pt="10px">{progress + "%"}</Center>
          </Box>
      </Box>
    </AspectRatio>
  );
}

const ProjectTab: React.FC<ProjectTabProps> = (props) => {

  const { project: currentProject, teamSlug } = props;

  const [project, setProject] = useState<Project>(currentProject);
  const [transformationIdx, setTransformationIdx] = useState(0);

  const transformations = project && project?.transformations;
  const transformation = transformations && transformations[transformationIdx];
  const parseTranscript = transformation && transformation.transcript && JSON.parse(transformation.transcript)

  const isProcessing = transformations?.length === 0 || (transformations[transformationIdx].status && transformations[transformationIdx].status !== "complete")

  const [getProjectData, { data, refetch }]
    = useLazyQuery(GET_CURRENT_PROJECT, { variables: { teamSlug, projectId: project?.id }, pollInterval: !isProcessing ? 0 : 10000, fetchPolicy: 'no-cache' });

  const [deleteTransformation, { loading: deleteTfnLoading } ] = useMutation(DELETE_TRANSFORMATION);

  useEffect(() => {
    const newProjectData = data?.getTeamById.projects[0];
    const newTransformations = newProjectData?.transformations;
    if (newTransformations?.length) {
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
            { (isProcessing && transformation) ? <LoadingBox progress={transformation.progress} /> : <VideoPlayer src={transformation ? transformation?.targetMedia : project.sourceMedia } onTimeUpdate={onTimeUpdate} /> }
            <HStack display="flex" flexWrap={"wrap"} overflow={"auto"} spacing={"10px"} pt="10px" hidden={!transformations?.length}>
              { transformations?.map((t: any, idx: number) => (
                <Button
                  key={idx}
                  onClick={() => setTransformationIdx(idx)}
                  variant={idx === transformationIdx ? "solid" : "outline"}
                  pointerEvents={idx === transformationIdx ? "none" : "auto"}
                >
                  {t?.targetLanguage}
                </Button>
              ))}
              <NewTransformationModel project={project} refetch={refetch} />
              <SingleActionModal
                heading={"Delete Dubbing"}
                body={`Are you sure you want to delete this dubbing? This action is irreversible.`}
                action={() => deleteDubbing(transformations[transformationIdx].id)}
                loading={deleteTfnLoading}
                isOpen={isOpen}
                onClose={onClose}
              />
              {!transformations[transformationIdx]?.isSource && <Button onClick={onOpen} leftIcon={<TrashIcon />} variant={"outline"}>Delete Dubbing</Button>}
            </HStack>
          </GridItem>

          <GridItem area={'transcript'} h="full" w="full" borderWidth={"1px"} rounded="lg" maxH={"596px"}>

          { !parseTranscript ?
            <LoadingTranscriptView />
            :
            <TranscriptView segments={parseTranscript?.segments} />
          }
          </GridItem>

        </Grid>
      </Box>
    </Box>
  );
};

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


const LoadingTranscriptView = () => {
  return (
    <VStack h="full" p="10px">
      <VStack overflow={"scroll"} h="full" position={"relative"} w="full">
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

const GET_TEAMS = gql`
  query GetTeams {
    getTeams {
      slug
      name
      projects {
        id
        title
      }
    }
  }
`;


interface ProjectDashboardProps {
  teamSlug: string;
  projectId: number;
}

const ProjectDashboard: NextPage<ProjectDashboardProps> = ({ teamSlug, projectId }) => {

  const { data: currentTeamsData, loading } = useQuery(GET_TEAMS);

  const { data: currentProjectData }
  = useQuery(GET_CURRENT_PROJECT, { variables: { teamSlug, projectId }, fetchPolicy: 'cache-and-network' });

  const currentProject: Project = currentProjectData?.getTeamById?.projects?.[0];

  const textColor = useColorModeValue("black", "white");
  const bgColor = useColorModeValue("white", "black");

  const teams = currentTeamsData?.getTeams;
  const projects = currentTeamsData?.getTeams.find((team: Team) => team.slug === teamSlug)?.projects;

  return (
    <Box>
      <Head>
        <title>Project | PlanetCast</title>
        <meta
          name="description"
          content="Cast Content in any Language, Across the Planet"
        />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Box position={"fixed"} top={0} left={0} w="full" p="10px" backgroundColor={useColorModeValue("white", "black")} zIndex={1000}>
        <Navbar projects={projects} teams={teams} teamSlug={teamSlug} projectId={projectId} />
      </Box>

      { !loading && currentTeamsData ?

        <Box pt={"70px"}>
          <Tabs
            variant="line"
            colorScheme="gray"
            isLazy
            onChange={(index) => {
              switch (index) {
              }
            }}
          >
            <Box
              position={'fixed'}
              w={'100%'}
              zIndex={'999'}
              display={"flex"}
              alignItems={"center"}
              justifyContent={"center"}
            >
              <TabList pl={'25px'} w="full" maxW={"1920px"} backgroundColor={bgColor}>
                <Tab><Text textColor={textColor}>Project</Text></Tab>
                <Tab><Text textColor={textColor}>Settings</Text></Tab>
              </TabList>
            </Box>
            <TabPanels overflow={'auto'} pt={10}>
              <TabPanel>
                { currentProject && <ProjectTab project={currentProject} teamSlug={teamSlug} /> }
              </TabPanel>
              <TabPanel>
                <SettingsTab projectId={projectId} teamSlug={teamSlug} />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </Box>

        :

        <Center pt="100px">
          <Spinner size={"lg"} />
        </Center>

      }

    </Box>
  );
};

export default ProjectDashboard;

export const getServerSideProps: GetServerSideProps= async ({ params }) => {

  const teamSlug = params?.teamSlug;
  const projectId = params?.projectId?.[0];

  if (projectId === undefined) {
    return {
      notFound: true
    }
  }

  return {
    props: {
      teamSlug,
      projectId: parseInt(projectId)
    }
  }
};
