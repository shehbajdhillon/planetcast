import { GetApolloClient } from "@/apollo-client";
import { GetServerSideProps, NextPage } from "next";

import { getAuth } from "@clerk/nextjs/server";
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
  mutation CreateTeam($teamType: TeamType!, $addTrial: Boolean!) {
    createTeam(teamType: $teamType, addTrial: $addTrial) {
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
    const { data } = await apolloClient.mutate({
      mutation: CREATE_TEAM,
      variables: {
        teamType: 'PERSONAL',
        addTrial: true,
      }
    });
    teams = [data.createTeam];
  }

  return {
    redirect: {
      permanent: false,
      destination: `dashboard/${teams[0].slug}`
    },
    props: { teams }
  }
};
