import { GetApolloClient } from "@/apollo-client";
import { GetServerSideProps, NextPage } from "next";

import { getAuth } from "@clerk/nextjs/server";
import { gql } from "@apollo/client";
import { useRouter } from "next/router";
import { useEffect } from "react";

interface PageProps {
  redirect: string;
};

const Index: NextPage<PageProps> = ({ redirect }) => {

  const router = useRouter();

  useEffect(() => {
    router.push(redirect);
  }, []);

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
    props: { redirect: `/dashboard/${teams[0].slug}` }
  }
};
