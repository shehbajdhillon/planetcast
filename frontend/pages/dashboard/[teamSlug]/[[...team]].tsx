import {
  Box,
  HStack,
  Skeleton,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  useColorModeValue,
} from "@chakra-ui/react";
import { NextPage, GetServerSideProps } from "next";

import DashboardTab from "@/components/dashboard/dashboard_tab";
import Head from "next/head";
import Navbar from "@/components/dashboard/navbar";
import { Project, Team } from "@/types";
import { gql, useQuery } from "@apollo/client";
import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import SettingsTab from "@/components/dashboard/settings_tab";

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

const GET_TEAM_BY_ID = gql`
  query GetTeamById($teamSlug: String!) {
    getTeamById(teamSlug: $teamSlug) {
      slug
      name
      projects {
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
      subscriptionPlans {
        id
        remainingCredits
        stripeSubscriptionId
        subscriptionData {
          currentPeriodStart
          currentPeriodEnd
          status
          interval
          planName
          costInUsd
          lastFourCardDigits
        }
      }
    }
  }
`;

export interface DashboardPageProps {
  teamSlug: string;
  tab: number;
};

const Dashboard: NextPage<DashboardPageProps> = ({ teamSlug, tab }) => {

  const { data, refetch: getAllTeamsRefetch } = useQuery(GET_TEAMS, { fetchPolicy: 'no-cache' });
  const {
    data: currentTeamData,
    refetch: currentTeamRefetch,
    loading,
    error
  } = useQuery(GET_TEAM_BY_ID, { fetchPolicy: 'no-cache', variables: { teamSlug } });

  const teams = data?.getTeams;
  const projects = data?.getTeams.find((team: Team) => team.slug === teamSlug)?.projects;

  const currentTeamProjects: Project[] = currentTeamData?.getTeamById?.projects;
  const currentTeam: Team = currentTeamData?.getTeamById;

  const [tabIdx, setTabIdx] = useState(tab);
  const router = useRouter();

  const textColor = useColorModeValue("black", "white");
  const bgColor = useColorModeValue("white", "black");

  const refetch = async () => {
    await currentTeamRefetch();
    await getAllTeamsRefetch();
  };

  useEffect(() => {
    if (error && error.message === "Access Denied") router.push('/404');
  }, [error]);

  useEffect(() => {
    const params = router.query.team;
    if (!params) {
      setTabIdx(0);
      return;
    }

    const [tabName] = params as string[];
    const tabProps = getTabIdx(tabName);

    if (!tabProps || tabProps.notFound) return;
    setTabIdx(tabProps.defaultTab);
  }, [router.query]);

  return (
    <Box>
      <Head>
        <title>Dashboard | PlanetCast</title>
        <meta
          name="description"
          content="Cast Content in any Language, Across the Planet"
        />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Box position={"fixed"} top={0} left={0} w="full" p="10px" backgroundColor={useColorModeValue("white", "black")} zIndex={1000}>
        <Navbar teamSlug={teamSlug} projects={projects} teams={teams} />
      </Box>

      <Box pt="80px">

        <Tabs
          variant="enclosed"
          colorScheme="gray"
          isLazy
          defaultIndex={tab}
          index={tabIdx}
          onChange={(index) => {
            switch (index) {
              case 0:
                router.push(`/dashboard/${teamSlug}`, undefined, { shallow: true });
                setTabIdx(index);
                break;
              case 1:
                router.push(`/dashboard/${teamSlug}/settings`, undefined, { shallow: true });
                setTabIdx(index);
                break;
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

              <Tab hidden={loading}><Text textColor={textColor}>Projects</Text></Tab>
              <Tab hidden={loading}><Text textColor={textColor}>Team Settings</Text></Tab>

              <HStack hidden={!loading} py="10px" spacing="15px">
                <Skeleton h="42px" w="100px" rounded={"md"} />
                <Skeleton h="42px" w="100px" rounded={"md"} />
              </HStack>

            </TabList>
          </Box>

          <TabPanels overflow={'auto'} pt={10}>
            <TabPanel>
              <DashboardTab teamSlug={teamSlug} projects={currentTeamProjects} refetch={refetch} />
            </TabPanel>
            <TabPanel>
              <SettingsTab
                teamSlug={teamSlug}
                params={router.query.team as string[]}
                refetch={refetch}
                loading={loading}
                team={currentTeam}
              />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>

    </Box>

  );
};

export default Dashboard;

const getTabIdx = (tabName: string) => {

  let defaultTab: number = 0;
  let notFound: boolean = false;

  switch(tabName) {
    case undefined:
      defaultTab = 0;
      break;
    case 'settings':
      defaultTab = 1;
      break;
    default:
      notFound = true;
      break;
  }

  return {
    defaultTab: defaultTab,
    notFound: notFound
  } as const;
};

export const getServerSideProps: GetServerSideProps = async ({ params }) => {

  const teamSlug = params?.teamSlug;
  const team = params?.team;

  let initialTabIdx = 0;

  if (team && team.length > 0) {
    const [tabName] = team as string[];
    const { defaultTab, notFound } = getTabIdx(tabName);
    if (notFound) return { notFound: true }
    initialTabIdx = defaultTab;
  }

  return {
    props: {
      teamSlug,
      tab: initialTabIdx
    }
  }
};

