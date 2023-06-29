import Navbar from "@/components/dashboard/navbar";
import { Project, Segment, Team } from "@/types";
import { gql, useMutation, useQuery } from "@apollo/client";
import {
  Box,
  Button,
  Center,
  Grid,
  GridItem,
  Heading,
  Spinner,
  Text,
  VStack,
  useColorModeValue,
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
import useWindowDimensions from "@/hooks/useWindowDimensions";


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
          transcript
        }
      }
    }
  }
`;

const ProjectTab: React.FC<ProjectTabProps> = ({ project, teamSlug }) => {

  const { data }
    = useQuery(GET_TRANSCRIPT, { variables: { teamSlug, projectId: project?.id } });

  const transcript = data?.getTeamById?.projects?.[0].transformations?.[0]?.transcript
  const parseTranscript = transcript && JSON.parse(transcript)

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

          <GridItem area={'video'} h="full" w="full" rounded={"lg"}>
            <VideoPlayer src={project?.sourceMedia} />
          </GridItem>

          <GridItem area={'transcript'} h="full" w="full" borderWidth={"1px"} rounded="lg" maxH={"596px"}>

          { !parseTranscript ?
            <Center h="full">
              <Heading>Fetching Transcript <Spinner /></Heading>
            </Center>
            :
            <VStack overflow={"auto"} p="10px" h="full">
              {parseTranscript?.segments.map((segment: Segment, idx: number) => (
                <Text
                  w="full"
                  borderWidth={"1px"}
                  rounded="lg"
                  p="10px"
                  key={idx}
                >
                  { segment.text.trim() }
                </Text>
              ))}
            </VStack>
          }
          </GridItem>

        </Grid>
      </Box>
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
          <Spinner />
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
