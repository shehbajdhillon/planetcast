import { GetApolloClient } from "@/apollo-client";
import { GetServerSideProps, NextPage } from "next";

import { v4 } from 'uuid';
import { clerkClient, getAuth } from "@clerk/nextjs/server";
import { gql } from "@apollo/client";

const Index: NextPage = () => {
  return <div />;
};

const GET_TEAMS = gql`
  query GetTeams {
    getTeams {
      slug
    }
  }
`;

const CREATE_TEAM = gql`
  mutation CreateTeam($name: String!, $slug: String!, $teamType: TeamType!) {
    createTeam(slug: $slug, name: $name, teamType: $teamType) {
      slug
    }
  }
`;

export default Index;

export const getServerSideProps: GetServerSideProps = async (ctx) => {

  const { getToken, userId } = getAuth(ctx.req)
  const apolloClient = GetApolloClient(true, getToken);

  let teams: any[] = [];

  const { data } = await apolloClient.query({ query: GET_TEAMS });
  teams = data.getTeams;

  if (userId && data?.getTeams?.length === 0) {
    const user = await clerkClient.users.getUser(userId);
    const { data } = await apolloClient.mutate({
      mutation: CREATE_TEAM,
      variables: {
        slug: v4(),
        name: `${user.firstName}'s Personal Workspace`,
        teamType: 'PERSONAL',
      }
    });
    teams = [data.createTeam];
  }

  return {
    redirect: {
      permanent: false,
      destination: `${teams[0].slug}`
    },
    props: { teams }
  }
};
