import Navbar from "@/components/dashboard/navbar";
import { Team } from "@/types";
import { gql, useQuery } from "@apollo/client";
import {
  Box,
  Heading,
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


const SettingsTab: React.FC = () => {
  return (
    <Box
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      w={"full"}
    >
      <Box w="full" maxW={"1920px"}>
        <Heading>SETTINGS</Heading>
      </Box>
    </Box>
  );
};

const ProjectTab: React.FC = () => {
  return (
    <Box
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      w={"full"}
    >
      <Box w="full" maxW={"1920px"}>
        <Heading>PROJECT</Heading>
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
      }
    }
  }
`;

interface ProjectDashboardProps {
  teamSlug: string;
  projectId: number;
}

const ProjectDashboard: NextPage<ProjectDashboardProps> = ({ teamSlug, projectId }) => {

  const { data } = useQuery(GET_TEAMS);

  const textColor = useColorModeValue("black", "white");

  const teams = data?.getTeams;
  const projects = data?.getTeams.find((team: Team) => team.slug === teamSlug)?.projects;

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
            <TabList pl={'25px'} w="full" maxW={"1920px"}>
              <Tab><Text textColor={textColor}>Project</Text></Tab>
              <Tab><Text textColor={textColor}>Settings</Text></Tab>
            </TabList>
          </Box>
          <TabPanels overflow={'auto'} pt={10}>
            <TabPanel>
              <ProjectTab />
            </TabPanel>
            <TabPanel>
              <SettingsTab />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>

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
