import Navbar from "@/components/dashboard/navbar";
import { Team } from "@/types";
import { gql, useQuery } from "@apollo/client";
import {
  Box,
  useColorModeValue,
} from "@chakra-ui/react";
import { GetServerSideProps, NextPage } from "next";

import Head from "next/head";

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
      <Box
        display={"flex"}
        justifyContent={"center"}
      >
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
