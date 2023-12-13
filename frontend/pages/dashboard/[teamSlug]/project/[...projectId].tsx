import Navbar from "@/components/dashboard/navbar";
import { Project, Team } from "@/types";
import { gql, useQuery } from "@apollo/client";
import {
  AspectRatio,
  Box,
  Center,
  Grid,
  GridItem,
  Skeleton,
  Text,
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
import ProjectTab, { LoadingTranscriptView } from "@/components/dashboard/project/project_tab";
import ProjectSettingsTab from "@/components/dashboard/project/project_settings_tab";
import { useEffect } from "react";
import { useRouter } from "next/router";

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

const GET_TEAMS = gql`
  query GetTeams {
    getTeams {
      slug
      name
      projects {
        id
        title
      }
      subscriptionPlans {
        remainingCredits
      }
    }
  }
`;


interface ProjectDashboardProps {
  teamSlug: string;
  projectId: number;
}

const ProjectDashboard: NextPage<ProjectDashboardProps> = ({ teamSlug, projectId }) => {

  const {
    data: allTeamsData,
    refetch: allTeamsRefetch,
  } = useQuery(GET_TEAMS);

  const {
    data: currentProjectData,
    loading: currentProjectLoading,
    refetch: currentProjectRefetch,
    error
  } = useQuery(GET_CURRENT_PROJECT, { variables: { teamSlug, projectId } });

  const refetch = async () => {
    await allTeamsRefetch();
    await currentProjectRefetch();
  };

  const currentProject: Project = currentProjectData?.getTeamById?.projects?.[0];
  const textColor = useColorModeValue("black", "white");
  const bgColor = useColorModeValue("white", "black");

  const allTeams = allTeamsData?.getTeams;
  const allProjects = allTeamsData?.getTeams.find((team: Team) => team.slug === teamSlug)?.projects;

  const router = useRouter();

  useEffect(() => {
    if (error && error.message === "Access Denied") router.push('/404');
  }, [error]);

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
        <Navbar projects={allProjects} teams={allTeams} teamSlug={teamSlug} projectId={projectId} />
      </Box>

      { (!currentProjectLoading && currentProjectData) ?

        <Box pt={"80px"}>
          <Tabs
            variant="enclosed"
            colorScheme="gray"
            isLazy
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
                <ProjectSettingsTab projectId={projectId} teamSlug={teamSlug} refetch={refetch} />
              </TabPanel>
            </TabPanels>
          </Tabs>
        </Box>
        :
        <Center pt="100px">
          <LoadingTabs />
        </Center>
      }
    </Box>
  );
};

const LoadingTabs = () => {
  return (
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
        p="1rem"
      >
        <GridItem area={'video'} h="full" w="full" rounded={"lg"} maxW={"1280px"}>
          <AspectRatio ratio={16/9}>
            <Skeleton w="full" h="100px" rounded={"lg"} />
          </AspectRatio>
        </GridItem>
        <GridItem area={'transcript'} h="full" w="full" borderWidth={"1px"} rounded="lg" maxH={"596px"}>
          <LoadingTranscriptView transcribing={false} />
        </GridItem>
      </Grid>
    </Box>
  );
};

export default ProjectDashboard;

export const getServerSideProps: GetServerSideProps = async ({ params }) => {

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
