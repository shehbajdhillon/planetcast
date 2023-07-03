import Navbar from "@/components/dashboard/navbar";
import { Project, Segment, Team } from "@/types";
import { gql, useMutation, useQuery } from "@apollo/client";
import {
  Box,
  Button,
  Center,
  Checkbox,
  Grid,
  GridItem,
  HStack,
  Heading,
  IconButton,
  Spinner,
  Text,
  VStack,
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
import { useEffect, useRef, useState } from "react";
import { PlusIcon } from "lucide-react";
import { useVideoSeekStore } from "@/stores/video_seek_store";
import { formatTime } from "@/utils";
import { text } from "stream/consumers";
import NewTransformationModel from "@/components/new_transformation_modal";


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

  const [deleteProjectMutation, {}] = useMutation(DELETE_PROJECT);

  const router = useRouter();

  const deleteProject = async () => {
    const res = await deleteProjectMutation({ variables: { projectId } });
    if (res) router.push(`/${teamSlug}`);
  };

  return (
    <Box
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      w={"full"}
    >
      <Box w="full" maxW={"1920px"}>
        <Center>
          <Button colorScheme="red" onClick={deleteProject}>
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

const GET_TRANSCRIPT = gql`
  query GetTranscript($teamSlug: String!, $projectId: Int64!) {
    getTeamById(teamSlug: $teamSlug) {
      projects(projectId: $projectId) {
        transformations {
          id
          targetMedia
          targetLanguage
          transcript
          isSource
        }
      }
    }
  }
`;

const ProjectTab: React.FC<ProjectTabProps> = ({ project, teamSlug }) => {

  const [currentTransformation, setCurrentTransformation] = useState(0);
  const [transformationPresent, setTranformationPresent] = useState(false);

  const { data, loading }
    = useQuery(GET_TRANSCRIPT, { variables: { teamSlug, projectId: project?.id }, pollInterval: transformationPresent ? 0 : 10000 });

  const currentProject: Project = data?.getTeamById?.projects?.[0];
  const transformationsArray = currentProject && currentProject.transformations;
  const transformation = transformationsArray && transformationsArray[currentTransformation];
  const parseTranscript = transformation && transformation.transcript && JSON.parse(transformation.transcript)

  const setCurrentSeek = useVideoSeekStore((state) => state.setCurrentSeek);
  const onTimeUpdate = (time: number) => setCurrentSeek(time);

  useEffect(() => {
    setCurrentSeek(0);
  }, [setCurrentSeek]);

  useEffect(() => {
    if (transformationsArray?.length) {
      setTranformationPresent(true);
    }
  }, [transformationsArray]);

  if (loading) {
    return (
      <Center h="full" w="full">
        <Spinner size={"lg"} />
      </Center>
    )
  }

  return (
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
            <VideoPlayer src={transformation ? transformation?.targetMedia : project.sourceMedia } onTimeUpdate={onTimeUpdate} />
            <HStack overflow={"auto"} spacing={"10px"} pt="10px" hidden={!transformationsArray?.length}>
              { transformationsArray?.map((t: any, idx: number) => (
                <Button
                  key={idx}
                  onClick={() => setCurrentTransformation(idx)}
                  variant={idx === currentTransformation ? "solid" : "outline"}
                  pointerEvents={idx === currentTransformation ? "none" : "auto"}
                >
                  {t?.targetLanguage}
                </Button>
              ))}
              <NewTransformationModel project={currentProject} />
            </HStack>
          </GridItem>

          <GridItem area={'transcript'} h="full" w="full" borderWidth={"1px"} rounded="lg" maxH={"596px"}>

          { !parseTranscript ?
            <Center h="full">
              <Heading>Generating Transcript <Spinner /></Heading>
            </Center>
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

  const bgColorHighlight = useColorModeValue("black", "white");
  const textColor = useColorModeValue("white", "black");

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
      <VStack overflow={"scroll"} h="full" id={parentId} position={"relative"}>
        {segments.map((segment: Segment, idx: number) => (
          <Button
            key={idx}
            rounded="10px"
            whiteSpace={'normal'}
            height="auto"
            blockSize={'auto'}
            w="full"
            justifyContent="left"
            leftIcon={<Text>{formatTime(segment.start)}</Text>}
            variant={highlight(segment) ? 'solid' : 'outline'}
            textColor={highlight(segment) ? textColor : 'inherit'}
            bgColor={highlight(segment) ? bgColorHighlight : 'inherit'}
            id={transcriptMessageId(segment)}
          >
            <Text
              key={idx}
              textAlign={"left"}
              padding={2}
            >
              { segment.text.trim() }
            </Text>
          </Button>
        ))}
      </VStack>
    </VStack>
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
        sourceMedia
        sourceLanguage
      }
    }
  }
`;

interface ProjectDashboardProps {
  teamSlug: string;
  projectId: number;
}

const ProjectDashboard: NextPage<ProjectDashboardProps> = ({ teamSlug, projectId }) => {

  const { data, loading } = useQuery(GET_TEAMS);

  const textColor = useColorModeValue("black", "white");
  const bgColor = useColorModeValue("white", "black");

  const teams = data?.getTeams;
  const projects = data?.getTeams.find((team: Team) => team.slug === teamSlug)?.projects;
  const currentProject = projects?.find((project: Project) => project.id === projectId);

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

      { !loading && data ?

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
